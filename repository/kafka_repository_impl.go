package repository

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaRepositoryImpl struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewKafkaRepository(brokers []string, topic, groupID string) KafkaRepository {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
	})

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:          brokers,
		Topic:            topic,
		GroupID:          groupID,
		ReadBatchTimeout: 10 * time.Second,
	})

	return &KafkaRepositoryImpl{
		writer: writer,
		reader: reader,
	}
}

func (k *KafkaRepositoryImpl) Publish(ctx context.Context, message string) error {
	return k.writer.WriteMessages(ctx, kafka.Message{
		Value: []byte(message),
	})
}

func (k *KafkaRepositoryImpl) Consume(ctx context.Context) (<-chan string, error) {
	messages := make(chan string)

	startTime := time.Now()
	go func() {
		defer close(messages)
		for {
			// Read a message from Kafka
			msg, err := k.reader.ReadMessage(ctx)
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}
			// Send the message to the messages channel
			messages <- string(msg.Value)
		}
	}()
	duration := time.Since(startTime)
	log.Printf("Execution consumer time: %v\n", duration)
	return messages, nil
}

func (k *KafkaRepositoryImpl) Close() error {
	if err := k.writer.Close(); err != nil {
		return err
	}
	return k.reader.Close()
}
