package service

import "context"

type KafkaService interface {
	SendMessage(ctx context.Context, message string) error
	ReceiveMessages(ctx context.Context)
}
