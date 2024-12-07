package repository

import "context"

type KafkaRepository interface {
	Publish(ctx context.Context, message string) error
	Consume(ctx context.Context) (<-chan string, error)
	Close() error
}
