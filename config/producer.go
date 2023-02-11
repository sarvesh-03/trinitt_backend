package config

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

var producer *kafka.Conn


func CreateProducer() {
	topic := "qwerty"
	partition := 0
	var err error
	producer, err = kafka.DialLeader(context.Background(), "tcp", "kafka:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
}

func GetProducer() *kafka.Conn{
	return producer
}
