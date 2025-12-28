package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	}

	p, err := kafka.NewProducer(config)
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					log.Println("Delivery failed:", m.TopicPartition)
				} else {
					log.Printf(
						"Delivered to %s [partition %d] @ offset %v | value=%s\n",
						*ev.TopicPartition.Topic,
						ev.TopicPartition.Partition,
						ev.TopicPartition.Offset,
						string(ev.Value),
					)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "helloworld"
	for n := 0; n < 100; n++ {
		fmt.Printf("send message to topic %s index %d", topic, n)
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte(strconv.Itoa(n)),
			Value:          []byte(fmt.Sprintf("Hello 2, Kafka%d :)", n)),
		}, nil)
	}

	// Wait for message deliveries
	p.Flush(5 * 1000)
}
