package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	startTime := time.Now()
	// Loop to send 10 messages
	for test := 0; test < 10000; test++ {
		message := fmt.Sprintf("Hello, Kafka with Repository Pattern! %d", test)

		// Print the message
		log.Println("Sending message:", message)

		// Send the message using the SendMessage method from kafkaService
		err := kafkaService.SendMessage(ctx, message)
		if err != nil {
			log.Println("Error sending message:", err)
			// Optionally, you can continue to the next iteration instead of returning
			continue
		}

		// // Send the message asynchronously using goroutine
		// go func(msg string) {
		// 	err := kafkaService.SendMessage(ctx, msg)
		// 	if err != nil {
		// 		log.Println("Error sending message:", err)
		// 	}
		// }(message)
	}
	duration := time.Since(startTime)
	log.Printf("Execution producer time: %v\n", duration)

	// Receive messages
	go kafkaService.ReceiveMessages(ctx)

	// if context is done or context is canceled
	<-ctx.Done()
	log.Println("Shutting down gracefully...")
}
