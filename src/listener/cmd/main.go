package main

import (
	"context"
	"frcofilippi/pedimeapp/listener/internal/handlers"
	"frcofilippi/pedimeapp/shared/config"
	"frcofilippi/pedimeapp/shared/events"
	"frcofilippi/pedimeapp/shared/logger"

	"go.uber.org/zap"
)

func main() {
	listener_logger := logger.InitLogger("listener-service")
	defer listener_logger.Sync()

	appConfig := config.NewApiConfiguration()

	dispatcher := events.Dispatcher{
		"ProductCreated": func(ctx context.Context, msg events.OutboxMessage) error {
			h := handlers.ProductCreatedHandler{}
			return h.HandleContext(ctx, msg)
		},
	}

	clientCfg := events.RabbitMQConfig{
		URL:        appConfig.Rabbitmqconfig.ConnectionStr,
		Exchange:   appConfig.Rabbitmqconfig.ExchangeName,
		DLExchange: appConfig.Rabbitmqconfig.DLExchange,
		DLQueue:    appConfig.Rabbitmqconfig.DLQueue,
	}
	eventClient, err := events.NewRabbitMQClient(clientCfg)
	if err != nil {
		listener_logger.Panic("consumer error", zap.String("error", err.Error()))
	}
	defer eventClient.Close()

	err = eventClient.Listen(context.Background(), dispatcher)
	if err != nil {
		listener_logger.Panic("listener error", zap.String("error", err.Error()))
	}

	listener_logger.Info("Listening for messages from broker")
}
