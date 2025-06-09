package main

import (
	"context"
	"frcofilippi/pedimeapp/listener/internal"
	"frcofilippi/pedimeapp/listener/internal/handlers"
	"frcofilippi/pedimeapp/shared/config"
	"frcofilippi/pedimeapp/shared/logger"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func main() {

	listener_logger := logger.InitLogger("listener-service")
	defer listener_logger.Sync()

	appConfig := config.NewApiConfiguration()

	hanlder := &handlers.ProductCreatedHandler{}

	con, err := connectToRabbitMq(appConfig.Rabbitmqconfig.ConnectionStr)
	if err != nil {
		listener_logger.Panic("rabbitmq error", zap.String("error", err.Error()))
	}
	defer con.Close()

	consumer, err := internal.NewConsumer(con, *appConfig.Rabbitmqconfig, hanlder)
	if err != nil {
		listener_logger.Panic("consumer error", zap.String("error", err.Error()))
	}

	err = consumer.ListenForMessages(context.Background())
	if err != nil {
		listener_logger.Panic("listener error", zap.String("error", err.Error()))
	}

	// log.Default().Printf("[Main] - Listening for messages")
	listener_logger.Info("Listening for messages from broker")
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
		logger.Error(
			"error connecting RabbitMQ",
			zap.Int("attempt", attempt),
			zap.Int("max_attempts", maxAttempts),
			zap.String("error", err.Error()),
		)

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
