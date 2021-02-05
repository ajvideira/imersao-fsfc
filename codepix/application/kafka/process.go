package kafka

import (
	"fmt"
	"log"
	"os"

	"github.com/ajvideira/imersao-fullstack-fullcycle/codepix/application/factory"
	appmodel "github.com/ajvideira/imersao-fullstack-fullcycle/codepix/application/model"
	"github.com/ajvideira/imersao-fullstack-fullcycle/codepix/domain/model"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jinzhu/gorm"
)

// KafkaProcessor represents a processor for kafka messages
type KafkaProcessor struct {
	Database *gorm.DB
	Producer *ckafka.Producer
	DeliveryChan chan ckafka.Event
}

// NewKafkaProcessor creates a new processor for kafka messages
func NewKafkaProcessor(database *gorm.DB, producer *ckafka.Producer, deliveryChan chan ckafka.Event) *KafkaProcessor {
	return &KafkaProcessor{
		Database: database,
		Producer: producer,
		DeliveryChan: deliveryChan,
	}
}

// Consume consumes messages
func (processor *KafkaProcessor) Consume() {
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers":os.Getenv("kafkaBootstrapServers"),
		"group.id":os.Getenv("kafkaConsumerGroupId"),
		"auto.offset.reset":"earliest",
	}
	consumer, err := ckafka.NewConsumer(configMap)
	if err != nil {
		panic(err)
	}

	topics := []string{os.Getenv("kafkaTransactionTopic"), os.Getenv("kafkaTransactionConfirmationTopic"), "teste"}
	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Kafka consumer has been started")
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			processor.processMessage(msg)
		}
	}
}

func (processor *KafkaProcessor) processMessage(msg *ckafka.Message) {
	transactionsTopic := os.Getenv("kafkaTransactionTopic")
	transactionsConfirmationTopic := os.Getenv("kafkaTransactionConfirmationTopic")

	switch topic := *msg.TopicPartition.Topic; topic {
	case transactionsTopic:
		processor.processTransaction(msg)
	case transactionsConfirmationTopic:
		processor.processTransactionConfirmation(msg)
	default:
		fmt.Println("unknown topic: ", topic)
		fmt.Println("message: ", string(msg.Value))
	}
}

func (processor *KafkaProcessor) processTransaction(msg *ckafka.Message) error {
	transaction := &appmodel.Transaction{};
	err := transaction.ParseJSON(msg.Value)
	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(processor.Database)
	createdTransaction, err := transactionUseCase.Register(
		transaction.AccountID,
		transaction.Amount,
		transaction.PixKeyTo,
		transaction.PixKeyKindTo,
		transaction.Description,
	)
	if err != nil {
		fmt.Println("error registering transaction", err)
		return err
	}

	topic := "bank" + createdTransaction.PixKeyTo.Account.Bank.Code
	transaction.ID = createdTransaction.ID
	transaction.Status = model.TransactionPending
	transactionJSONBytes, err := transaction.ToJSON();
	if err != nil {
		fmt.Println("error converting transaction to json", err)
	}

	err = Publish(string(transactionJSONBytes), topic, processor.Producer, processor.DeliveryChan)
	if err != nil {
		return err
	}
	return nil
}

func (processor *KafkaProcessor) processTransactionConfirmation(msg *ckafka.Message) error {
	transaction := &appmodel.Transaction{};
	err := transaction.ParseJSON(msg.Value)
	if err != nil {
		return err
	}

	if transaction.Status == model.TransactionConfirmed {
		err = processor.confirmTransaction(transaction)
		if err != nil {
			return err
		}
	} else if transaction.Status == model.TransactionCompleted {
		err = processor.completeTransaction(transaction)
	}

	return nil
}

func (processor *KafkaProcessor) confirmTransaction(transaction *appmodel.Transaction) error {
	transactionUseCase := factory.TransactionUseCaseFactory(processor.Database)
	
	confirmedTransaction, err := transactionUseCase.Confirm(transaction.ID)
	if err != nil {
		fmt.Println("error confirming transaction", err)
		return err
	}

	topic := "bank" + confirmedTransaction.AccountFrom.Bank.Code
	transaction.Status = model.TransactionConfirmed
	transactionJSONBytes, err := transaction.ToJSON();
	if err != nil {
		fmt.Println("error converting transaction to json", err)
	}
	err = Publish(string(transactionJSONBytes), topic, processor.Producer, processor.DeliveryChan)
	if err != nil {
		return err
	}
	return nil
}

func (processor *KafkaProcessor) completeTransaction(transaction *appmodel.Transaction) error {
	transactionUseCase := factory.TransactionUseCaseFactory(processor.Database)
	
	_, err := transactionUseCase.Complete(transaction.ID)
	if err != nil {
		fmt.Println("error confirming transaction", err)
		return err
	}
	return nil
}