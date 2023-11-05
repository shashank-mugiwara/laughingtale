package conf

import (
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func InitKafkaConsumer() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092,localhost:9093",
		"group.id":          "foo",
		"auto.offset.reset": "smallest"})

	if err != nil {
		log.Fatalln("Unable to get kafka connection to consume messages")
		log.Panic("Unable to get kafka connection to consume messages")
	}

	subscription_topics := []string{"walconnect1.inventory.customers"}
	err = consumer.SubscribeTopics(subscription_topics, nil)

	if err != nil {
		log.Fatalln("Unable to get kafka connection to consume messages")
		log.Panic("Unable to get kafka connection to consume messages")
	}

	log.Println("Connection to kafka node is successful. Starting to listen to messages")

	run := true

	for run == true {
		ev := consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			log.Println("Got message from kafka ...")
			log.Println(string(e.Value))
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			run = false
		default:
			fmt.Printf("Ignored %v\n", e)
		}
	}
}
