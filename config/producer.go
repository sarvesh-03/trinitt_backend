package config

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var KafkaProducer *kafka.Producer 


func CreateProducer() {
	
	var err error
	
	KafkaProducer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaUrl})
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	
}

func GetProducer() *kafka.Producer{
	return KafkaProducer
}
