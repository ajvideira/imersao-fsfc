package kafka

import (
	"fmt"
	"os"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

// NewKafkaProducer creates a new kafka producer
func NewKafkaProducer() *ckafka.Producer {
	producer, err := ckafka.NewProducer(&ckafka.ConfigMap{
		"bootstrap.servers":os.Getenv("kafkaBootstrapServers"),
	})
	if err != nil {
		panic(err)
	}
	return producer
}

// Publish publishes a message to a topic
func Publish (msg string, topic string, producer *ckafka.Producer, deliveryChan chan ckafka.Event) error {
	message := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{
			Topic: &topic,
			Partition: ckafka.PartitionAny,
		},
		Value: []byte(msg),
	}
	err := producer.Produce(message, deliveryChan)
	if err != nil {
		return err
	}
	return nil
}

// DeliveryReport monitories a channel kafka
func DeliveryReport(deliveryChan chan ckafka.Event) {
	for e := range deliveryChan {
		switch ev := e.(type) {
		case *ckafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Println("Delivery message failed: ", ev.TopicPartition)
			} else {
				fmt.Println("Delivered message to: ", ev.TopicPartition)
			}
		}
	}
}