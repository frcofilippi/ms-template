package main

import (
	"context"
	"frcofilippi/pedimeapp/internal/application"
	"frcofilippi/pedimeapp/internal/common"
	"frcofilippi/pedimeapp/internal/events"
	"frcofilippi/pedimeapp/internal/product"
	"log"
	"net/http"

	"github.com/joho/godotenv"
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

	eventConsumer, err := events.NewRabbitMqConnection(
		appConfig.rabbitmqconfig.connectionStr,
		appConfig.rabbitmqconfig.exchangeName,
		appConfig.rabbitmqconfig.queue,
	)
	if err != nil {
		log.Fatalf("Error creating event consumer: %s", err.Error())
	}
	defer eventConsumer.Close()

	dispatcher := &application.MessageDispatcher{}
	listener := application.NewMessageListener(eventConsumer, dispatcher)

	listnerErrChan := make(chan error, 1)
	go func() {
		err = listener.Start()
		if err != nil {
			listnerErrChan <- err
		}
	}()

	serverErrChan := make(chan error, 1)
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			serverErrChan <- err
		}
	}()

	select {
	case err := <-listnerErrChan:
		log.Fatalf("error: %s", err.Error())
	case err := <-serverErrChan:
		log.Fatalf("error: %s", err.Error())
	}
}
