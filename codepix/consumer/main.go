package consumer

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// ConsumeMsgs consumes msgs of a topic
func ConsumeMsgs(producer *kafka.Producer, deliveryChan chan kafka.Event) {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":"kafka:9092",
		"group.id":"codepix",
		"auto.offset.reset":"earliest",
	}

	consumer, err := kafka.NewConsumer(configMap)
	if err != nil {
		panic(err)
	}

	topics := []string{"desafio2"}

	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Kafka Consumer has been started")

	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			panic(err)
		}
		switch topic := *msg.TopicPartition.Topic; topic {
		case "desafio2":
			fmt.Println("Topic: desafio2")
			fmt.Println("Message: ",string(msg.Value))
		default:
			fmt.Println("Unknown topic: ", topic)
			fmt.Println("Message: ",string(msg.Value))
		}
	}
}