package controllers

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/echo/v4"

	"github.com/hamba/avro/v2"
	"github.com/trinitt/config"
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
	"name": "Record",
	"fields": [
	  {
		"name": "user_id",
		"type": "int"
	  },
	  {
		"name": "entity_id",
		"type": "int"
	  },
	  {
		"name": "param",
		"type": {
		  "type": "array",
		  "items": {
			"type": "record",
			"namespace": "Record",
			"name": "param",
			"fields": [
			  {
				"name": "data_type",
				"type": "string"
			  },
			  {
				"name": "value",
				"type": "string"
			  }
			]
		  }
		}
	  }
	]
  }`

type ParamType struct {
	Data_type string `avro:"data_type" json:"data_type"`
	Value     string `avro:"value" json:"value"`
}

type Record struct {
	User_id   int      `avro:"user_id" json:"user_id"`
	Entity_id int      `avro:"entity_id" json:"entity_id"`
	Param     []ParamType `avro:"param" json:"param"`
}


func Produce(c echo.Context) error{

	schema, err := avro.Parse(Schema)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
	
	fmt.Println(schema)
	in := Record{

		User_id:   1,
		Entity_id: 1,
		Param: []ParamType{
			{
				Data_type: "INT",
				Value:     "5",
			},
			{
				Data_type: "string",
				Value:     "hello",
			},
		},
	}

	data, err := avro.Marshal(schema, in)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", data)
	topic := "myTopic"

	config.GetProducer().Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: 0},
		Value:          data,
	}, nil)

	out := Record{}
	err = avro.Unmarshal(schema, data, &out)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.Entity_id)

	return c.JSON(200,"qwerty")
}


