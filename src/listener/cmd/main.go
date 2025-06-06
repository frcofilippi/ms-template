package main

import (
	"context"
	"frcofilippi/pedimeapp/listener/internal"
	"frcofilippi/pedimeapp/listener/internal/handlers"
	"frcofilippi/pedimeapp/shared/config"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	appConfig := config.NewApiConfiguration()

	// hanlder := &internal.CustomMessageHandler{}
	hanlder := &handlers.ProductCreatedHandler{}

	con, err := connectToRabbitMq(appConfig.Rabbitmqconfig.ConnectionStr)
	if err != nil {
		panic(err)
	}
	defer con.Close()

	consumer, err := internal.NewConsumer(con, *appConfig.Rabbitmqconfig, hanlder)
	if err != nil {
		panic(err)
	}

	err = consumer.ListenForMessages(context.Background())
	if err != nil {
		panic(err)
	}

	log.Default().Printf("[Main] - Listening for messages")
}

func connectToRabbitMq(conStr string) (*amqp.Connection, error) {
	const (
		maxAttempts    = 5
		initialBackoff = 1 * time.Second
		maxBackoff     = 16 * time.Second
	)
	var (
		con     *amqp.Connection
		err     error
		backoff = initialBackoff
	)
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		con, err = amqp.Dial(conStr)
		if err == nil {
			return con, nil
		}
		log.Printf("Failed to connect to RabbitMQ (attempt %d/%d): %v", attempt, maxAttempts, err)
		if attempt < maxAttempts {
			time.Sleep(backoff)
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}
	return nil, err
}
