package cmd

import (
	"os"

	"github.com/ajvideira/imersao-fullstack-fullcycle/codepix/application/grpc"
	"github.com/ajvideira/imersao-fullstack-fullcycle/codepix/application/kafka"
	"github.com/ajvideira/imersao-fullstack-fullcycle/codepix/infrastructure/db"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/cobra"
)

var gRPCportNumber int

// allCmd represents the all command
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Run gRPC and Kafka",
	Run: func(cmd *cobra.Command, args []string) {
		database := db.ConnectDB(os.Getenv("env"))
		go grpc.StartGrpcServer(database, gRPCportNumber)

		deliveryChan := make(chan ckafka.Event)
		go kafka.DeliveryReport(deliveryChan)
		producer := kafka.NewKafkaProducer();
		//kafka.Publish("Teste de funcionamento", "teste", producer, deliveryChan)

		kafkaProcessor := kafka.NewKafkaProcessor(database, producer, deliveryChan)
		kafkaProcessor.Consume()
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
	allCmd.Flags().IntVarP(&gRPCportNumber, "port", "p", 50051, "gRPC Server Port")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// allCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// allCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
