package main

import (
	"context"
	"frcofilippi/pedimeapp/internal/application"
	"frcofilippi/pedimeapp/internal/common"
	"frcofilippi/pedimeapp/internal/product"
	"frcofilippi/pedimeapp/shared/config"
	"frcofilippi/pedimeapp/shared/events"
	"log"
	"net/http"
)

func main() {

	appConfig := config.NewApiConfiguration()
	sqlConnection, err := common.NewPostgresConnection(appConfig.Dbconfig.ConnectionStr)
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
		Addr:    appConfig.Port,
		Handler: app.Mount(),
	}

	log.Printf("Configuration Loaded. Port: %s - Db: %s \n", appConfig.Port, appConfig.Dbconfig.ConnectionStr)
	log.Printf("Server running on port: %s \n", appConfig.Port)

	eventPublisher, err := events.NewRabbitMqConnection(
		appConfig.Rabbitmqconfig.ConnectionStr,
		appConfig.Rabbitmqconfig.ExchangeName,
	)

	if err != nil {
		log.Fatalf("Error creating event publisher: %s", err.Error())
	}
	defer eventPublisher.Close()

	obrelay := application.NewOutboxRelay(sqlConnection, eventPublisher)

	obrelay.Start(context.Background())

	serverErrChan := make(chan error, 1)
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			serverErrChan <- err
		}
	}()

	err = <-serverErrChan
	log.Fatalf("error: %s", err.Error())

}
