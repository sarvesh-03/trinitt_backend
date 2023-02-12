package controllers

import (
	"fmt"
	"time"

	"github.com/hamba/avro/v2"
	"github.com/trinitt/config"
	"github.com/trinitt/sockets"
)

func Consume() {
	schema, _ := avro.Parse(Schema1)

	go func() {
		run := true
		for run {
			msg, err := config.GetConsumer().ReadMessage(time.Second)
			if err == nil {
				fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
				out := Nodes{}
				avro.Unmarshal(schema, msg.Value, &out)
				fmt.Println(out)
				PostConsume(out)
			}
		}
	}()
}

func GetColorFromClusterID(id int) string {
	switch id {
	case 0:
		return "red"
	case 1:
		return "blue"
	case 2:
		return "green"
	case 3:
		return "yellow"
	case 4:
		return "purple"
	case 5:
		return "orange"
	case 6:
		return "pink"
	case 7:
		return "brown"
	case 8:
		return "grey"
	case 9:
		return "black"
	default:
		return "white"
	}
}

func PostConsume(rec Nodes) {
	nodes := []Node{}

	for index, node := range rec.Param {
		nodes = append(nodes, Node{
			ID:    index,
			X:     int(node.X),
			Y:     int(node.Y),
			Color: GetColorFromClusterID(node.Cluster),
		})
	}

	sockets.SendMessageToClient(1, "NODES", nodes)
	fmt.Println(rec)
}
