package controllers

import (
	"fmt"
	"time"

	"github.com/hamba/avro/v2"
	"github.com/trinitt/config"
)


func Consume(){
	schema, _ := avro.Parse(Schema)
	
	go func(){
	run:=true
	for run {
		msg, err := config.GetConsumer().ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			out := Record{}
			avro.Unmarshal(schema, msg.Value, &out)
			fmt.Println(out)
			PostConsume(out)
		}
	}
	}()
}

func PostConsume(rec Record){
	// Setup(rec)
	// clusters:=GetClustersForUser(1)
	// for _, cluster := range clusters {
	// 	fmt.Println("Cluster:")
	// 	for _, point := range cluster {
	// 		fmt.Println(point.(*dbscan.NamedPoint).Name)
	// 	}
	// }
	// Produce();
	fmt.Println(rec)
}