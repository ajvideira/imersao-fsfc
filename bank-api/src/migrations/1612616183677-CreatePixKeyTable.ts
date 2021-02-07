import { MigrationInterface, QueryRunner, Table, TableForeignKey } from 'typeorm';

export class CreatePixKeyTable1612616183677 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.createTable(
      new Table({
        name: 'pix_keys',
        columns: [
          {
            name: 'id',
            type: 'uuid',
            isPrimary: true,
          },
          {
            name: 'kind',
            type: 'varchar',
          },
          {
            name: 'key',
            type: 'varchar',
          },
          {
            name: 'bank_account_id',
            type: 'uuid',
          },
          {
            name: 'created_at',
            type: 'timestamp',
            default: 'CURRENT_TIMESTAMP',
          },
        ],
      }),
    );
    await queryRunner.createForeignKey(
      'pix_keys',
      new TableForeignKey({
        name: 'pix_keys_bank_account_id_foreign_key',
        columnNames: ['bank_account_id'],
        referencedTableName: 'bank_accounts',
        referencedColumnNames: ['id'],
      }),
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.dropTable('pix_keys');
  }
}
