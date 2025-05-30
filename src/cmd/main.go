package main

import (
	"context"
	"encoding/json"
	"frcofilippi/pedimeapp/internal/application"
	"frcofilippi/pedimeapp/internal/common"
	"frcofilippi/pedimeapp/internal/events"
	"frcofilippi/pedimeapp/internal/product"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appConfig := NewApiConfiguration()
	sqlConnection, err := common.NewPostgresConnection(appConfig.dbconfig.connectionStr)
	if err != nil {
		log.Fatalf("error connecting the database. Error %v", err)
	}

	productRepository, err := product.NewProductRepositoryWithCustomer(sqlConnection)
	if err != nil {
		log.Fatalf("Something went wrong when initializing the repository. Error: %v", err)
	}

	productService := product.NewProductService(productRepository)

	productRouter := product.NewProductRouter(productService)

	app := application.New(productRouter)

	server := &http.Server{
		Addr:    appConfig.port,
		Handler: app.Mount(),
	}

	log.Printf("Configuration Loaded. Port: %s - Db: %s \n", appConfig.port, appConfig.dbconfig.connectionStr)
	log.Printf("Server running on port: %s \n", appConfig.port)

	eventPublisher, err := events.NewRabbitMqConnection(
		appConfig.rabbitmqconfig.connectionStr,
		appConfig.rabbitmqconfig.exchangeName,
		appConfig.rabbitmqconfig.queue,
	)
	if err != nil {
		log.Fatalf("Error creating event publisher: %s", err.Error())
	}
	defer eventPublisher.Close()

	obrelay := application.NewOutboxRelay(sqlConnection, eventPublisher)

	obrelay.Start(context.Background())

	// consummer
	eventConsumer, err := events.NewRabbitMqConnection(
		appConfig.rabbitmqconfig.connectionStr,
		appConfig.rabbitmqconfig.exchangeName,
		appConfig.rabbitmqconfig.queue,
	)
	if err != nil {
		log.Fatalf("Error creating event consumer: %s", err.Error())
	}
	defer eventConsumer.Close()

	go func() {
		err = eventConsumer.Consume(func(delivery amqp.Delivery) {
			log.Default().Printf("Received message: %s", delivery.Body)

			var message application.OutboxMessage
			err := json.Unmarshal(delivery.Body, &message)
			if err != nil {
				log.Default().Printf("Error unmarshalling event: %s", err.Error())
				delivery.Acknowledger.Nack(delivery.DeliveryTag, false, false)
				return
			}
			log.Default().Printf("Event type: %s", message.EventType)
			var pCreated events.ProductCreatedEvent
			err = json.Unmarshal(message.Payload, &pCreated)
			if err != nil {
				log.Default().Printf("Error unmarshalling event: %s", err.Error())
				delivery.Acknowledger.Nack(delivery.DeliveryTag, false, false)
				return
			}
			log.Default().Printf("Message received and parsed!! %v \n", pCreated)
			delivery.Acknowledger.Ack(delivery.DeliveryTag, false)

		})
		if err != nil {
			log.Default().Fatalf("Error consuming messages: %s", err.Error())
		}
	}()

	err = server.ListenAndServe()

	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

}
