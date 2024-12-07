package repository

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaRepositoryImpl struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewKafkaRepository(brokers []string, topic, groupID string) KafkaRepository {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
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
	go func() {
		defer close(messages)
		for {
			msg, err := k.reader.ReadMessage(ctx)
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}
			messages <- string(msg.Value)
		}
	}()
	return messages, nil
}

func (k *KafkaRepositoryImpl) Close() error {
	if err := k.writer.Close(); err != nil {
		return err
	}
	return k.reader.Close()
}
