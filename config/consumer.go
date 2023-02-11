package config

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var KafkaConsumer *kafka.Consumer


func CreateConsumer() {
	
	var err error
	
	KafkaConsumer, err = kafka.NewConsumer(&kafka.ConfigMap{"bootstrap.servers": kafkaUrl,
	"group.id":          "myGroup",
	"auto.offset.reset": "earliest",})
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	KafkaConsumer.SubscribeTopics([]string{"yTopic", "^aRegex.*[Tt]opic"}, nil)
	
}

func GetConsumer() *kafka.Consumer{
	return KafkaConsumer
}

