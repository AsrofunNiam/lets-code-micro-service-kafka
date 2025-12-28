package main

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "helloworld_g1",
		"auto.offset.reset": "smallest",
	}

	c, err := kafka.NewConsumer(config)
	if err != nil {
		panic(err)
	}

	defer c.Close()

	c.SubscribeTopics([]string{"helloworld"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err != nil {
			panic(err)
		}

		println(
			"Topic:", *msg.TopicPartition.Topic,
			"| Partition:", msg.TopicPartition.Partition,
			"| Offset:", msg.TopicPartition.Offset,
			"| Value:", string(msg.Value),
		)
	}

}
