package producer

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

//NewKafkaProducer creates a new kafka producer
func NewKafkaProducer() *kafka.Producer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":"kafka:9092",
	})
	if err != nil {
		panic(err)
	}
	return producer;
}

//PublishMsg publishes a new message to a topic
func PublishMsg(msg string, topic string, producer *kafka.Producer, deliveryChan chan kafka.Event) error {
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value: []byte(msg),
 	}

	err := producer.Produce(message, deliveryChan)
	if err != nil {
		return err
	}
	return nil
}

// DeliveryWatch watches a delivery chanel for messages
func DeliveryWatch(deliveryChan chan kafka.Event) {
	for e := range deliveryChan{
		switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Println("Delivery message failed: ", ev.TopicPartition)
				} else {
					fmt.Println("Delivered message to: ", ev.TopicPartition)
				}
		}
	}
}