package service

import (
	"context"
	"log"

	"github.com/AsrofunNiam/lets-code-micro-service-kafka/repository"
)

type KafkaServiceImpl struct {
	KafkaRepository repository.KafkaRepository
}

func NewKafkaService(kafkaRepository repository.KafkaRepository) KafkaService {
	return &KafkaServiceImpl{KafkaRepository: kafkaRepository}
}

func (s *KafkaServiceImpl) SendMessage(ctx context.Context, message string) error {
	return s.KafkaRepository.Publish(ctx, message)
}

func (s *KafkaServiceImpl) ReceiveMessages(ctx context.Context) {
	messages, err := s.KafkaRepository.Consume(ctx)
	if err != nil {
		log.Println("Error consuming messages:", err)
		return
	}

	for msg := range messages {
		log.Println("Received message:", msg)
	}
}
