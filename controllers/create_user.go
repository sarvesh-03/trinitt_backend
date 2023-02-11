package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"
	"github.com/trinitt/config"
	"github.com/trinitt/utils"
)

type SignupRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	IsSeller  bool   `json:"is_seller"`
}

func SignupUser(c echo.Context) error {


	producer:= config.GetProducer()
	
	producer.SetWriteDeadline(time.Now().Add(10*time.Second))
	_, err := producer.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	// if err := producer.Close(); err != nil {
	// 	log.Fatal("failed to close writer:", err)
	// }

	return utils.SendResponse(c, http.StatusOK, "User created successfully")
}
