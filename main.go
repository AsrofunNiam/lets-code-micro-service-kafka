package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	config "github.com/AsrofunNiam/lets-code-micro-service-kafka/configuration"
	"github.com/AsrofunNiam/lets-code-micro-service-kafka/repository"
	"github.com/AsrofunNiam/lets-code-micro-service-kafka/service"
)

func main() {
	// Configuration
	kafkaConfig := config.NewKafkaConfig()

	// Set initialize
	repo := repository.NewKafkaRepository(kafkaConfig.Brokers, kafkaConfig.Topic, kafkaConfig.GroupID)
	defer repo.Close()

	kafkaService := service.NewKafkaService(repo)

	// Cancel context if implement goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()

	// Producer temp
	err := kafkaService.SendMessage(ctx, "Hello, Kafka with Repository Pattern!")
	if err != nil {
		log.Println("Error sending message:", err)
		return
	}

	// Consumer temp
	kafkaService.ReceiveMessages(ctx)

	// if context is done or context is canceled
	<-ctx.Done()
	log.Println("Shutting down gracefully...")
}
