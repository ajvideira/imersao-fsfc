package cmd

import (
	"time"

	"github.com/ajvideira/imersao-fullstack-fullcycle/codepix/consumer"
	"github.com/ajvideira/imersao-fullstack-fullcycle/codepix/producer"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/cobra"
)

// desafio2Cmd represents the desafio2 command
var desafio2Cmd = &cobra.Command{
	Use:   "desafio2",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		deliveriChan := make(chan kafka.Event)
		go producer.DeliveryWatch(deliveriChan);

		kafkaProducer := producer.NewKafkaProducer()
		go sendInfinityMsgs(kafkaProducer, deliveriChan)
		
		consumer.ConsumeMsgs(kafkaProducer, deliveriChan)
	},
}

func sendInfinityMsgs(kafkaProducer *kafka.Producer, deliveryChan chan kafka.Event) {
	for {
		time.Sleep(1 * time.Second)
		producer.PublishMsg("Echo!", "desafio2", kafkaProducer, deliveryChan)
	}
}

func init() {
	rootCmd.AddCommand(desafio2Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// desafio2Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// desafio2Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
