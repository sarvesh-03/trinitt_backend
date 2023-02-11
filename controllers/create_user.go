package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hamba/avro"
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

var Schema = `{
	"type": "record",
	"name": "simple",
	"namespace": "org.hamba.avro",
	"fields" : [
		{"name": "a", "type": "long"},
		{"name": "b", "type": "string"}
	]
}`

type SimpleRecord struct {
	A int64  `avro:"a"`
	B string `avro:"b"`
}


func SignupUser(c echo.Context) error {

	schema, err := avro.Parse(Schema)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
	producer:= config.GetProducer()
	
	in := SimpleRecord{A: 27, B: "foo"}

	data, err := avro.Marshal(schema, in)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", data)



	// Writing OCF data
	
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// // if err := producer.Close(); err != nil {
	// // 	log.Fatal("failed to close writer:", err)
	// // }
	// fmt.Println(ocfFileContents.String())
	producer.SetWriteDeadline(time.Now().Add(10*time.Second))
	_, err = producer.WriteMessages(
		kafka.Message{Value: data},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	return utils.SendResponse(c, http.StatusOK, "User created successfully")
}
