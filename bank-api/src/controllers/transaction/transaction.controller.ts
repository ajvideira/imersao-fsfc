import {
  Body,
  Controller,
  Get,
  Inject,
  OnModuleDestroy,
  OnModuleInit,
  Param,
  ParseUUIDPipe,
  Post,
  ValidationPipe,
} from '@nestjs/common';
import { ClientGrpc, ClientKafka, MessagePattern, Payload } from '@nestjs/microservices';
import { Producer } from '@nestjs/microservices/external/kafka.interface';
import { InjectRepository } from '@nestjs/typeorm';
import { TransactionDTO } from 'src/dto/transaction.dto';
import { BankAccount } from 'src/models/bank-account.model';
import { PixKey } from 'src/models/pix-key.model';
import { Transaction, TransactionOperation, TransactionStatus } from 'src/models/transaction.model';
import { Repository } from 'typeorm';

@Controller('/bank-accounts/:bankAccountId/transactions')
export class TransactionController implements OnModuleInit, OnModuleDestroy {
  private kafkaProducer: Producer;

  constructor(
    @InjectRepository(Transaction)
    private transactionRepo: Repository<Transaction>,
    @InjectRepository(BankAccount)
    private bankAccountRepo: Repository<BankAccount>,
    @InjectRepository(PixKey)
    private pixKeyRepo: Repository<PixKey>,
    @Inject('CODEPIX_PACKAGE')
    private client: ClientGrpc,
    @Inject('TRANSACTION_SERVICE')
    private kafkaClient: ClientKafka,
  ) {}

  async onModuleInit() {
    this.kafkaProducer = await this.kafkaClient.connect();
  }

  async onModuleDestroy() {
    await this.kafkaProducer.disconnect();
  }

  @Get()
  index(
    @Param('bankAccountId', new ParseUUIDPipe({ version: '4' }))
    bankAccountId: string,
  ) {
    return this.transactionRepo.find({
      where: {
        bank_account_id: bankAccountId,
      },
      order: {
        created_at: 'DESC',
      },
    });
  }
  @Post()
  async store(
    @Param('bankAccountId', new ParseUUIDPipe({ version: '4' }))
    bankAccountId: string,
    @Body(new ValidationPipe({ errorHttpStatusCode: 422 }))
    body: TransactionDTO,
  ) {
    await this.bankAccountRepo.findOneOrFail(bankAccountId);

    const transaction = this.transactionRepo.create({
      ...body,
      amount: body.amount * -1,
      bank_account_id: bankAccountId,
      operation: TransactionOperation.debit,
      status: TransactionStatus.pending,
    });
    await this.transactionRepo.save(transaction);
    console.log(transaction);

    const sendData = {
      id: transaction.external_id,
      accountId: bankAccountId,
      amount: body.amount,
      pixKeyTo: body.pix_key_key,
      pixKeyKindTo: body.pix_key_kind,
      description: body.description,
      status: transaction.status,
    };

    await this.kafkaProducer.send({
      topic: 'transactions',
      messages: [
        {
          key: 'transactions',
          value: JSON.stringify(sendData),
        },
      ],
    });

    return transaction;
  }

  @MessagePattern(`bank${process.env.BANK_CODE}`)
  async consumeTransactions(@Payload() message) {
    console.log('entrou em consumeTransactions()');
    console.log(message.value);
    if (message.value.status == TransactionStatus.pending) {
      await this.confirmTransaction(message.value);
    } else if (message.value.status == TransactionStatus.confirmed) {
      await this.completeTransaction(message.value);
    }
  }

  async confirmTransaction(data: any) {
    console.log('entrou em confirmTransaction()');
    const pixKey = await this.pixKeyRepo.findOneOrFail({
      where: {
        key: data.pixKeyTo,
        kind: data.pixKeyKindTo,
      },
    });

    const transaction = this.transactionRepo.create({
      external_id: data.id,
      amount: data.amount,
      description: data.description,
      bank_account_id: pixKey.bank_account_id,
      bank_account_from_id: data.accountId,
      operation: TransactionOperation.credit,
      status: TransactionStatus.completed,
    });

    await this.transactionRepo.save(transaction);

    const sendData = {
      ...data,
      status: TransactionStatus.confirmed,
    };

    await this.kafkaProducer.send({
      topic: 'transaction-confirmation',
      messages: [
        {
          key: 'transaction-confirmation',
          value: JSON.stringify(sendData),
        },
      ],
    });
  }

  async completeTransaction(data: any) {
    console.log('entrou em completeTransaction()');
    const transaction = await this.transactionRepo.findOneOrFail({
      where: {
        external_id: data.id,
        status: TransactionStatus.pending,
      },
    });

    await this.transactionRepo.update({ external_id: data.id }, { status: TransactionStatus.completed });

    const sendData = {
      id: data.id,
      accountId: transaction.bank_account_id,
      amount: Math.abs(transaction.amount),
      pixKeyTo: transaction.pix_key_key,
      pixKeyKindTo: transaction.pix_key_kind,
      description: transaction.description,
      status: TransactionStatus.completed,
    };

    await this.kafkaProducer.send({
      topic: 'transaction-confirmation',
      messages: [
        {
          key: 'transaction-confirmation',
          value: JSON.stringify(sendData),
        },
      ],
    });
  }
}
