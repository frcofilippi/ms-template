package main

import (
	"context"
	"frcofilippi/pedimeapp/internal/application"
	"frcofilippi/pedimeapp/internal/common"
	"frcofilippi/pedimeapp/internal/product"
	"frcofilippi/pedimeapp/shared/config"
	"frcofilippi/pedimeapp/shared/events"
	"frcofilippi/pedimeapp/shared/logger"
	"net/http"

	"go.uber.org/zap"
)

func main() {

	apilogger := logger.InitLogger("api-service")
	defer apilogger.Sync()

	appConfig := config.NewApiConfiguration()
	sqlConnection, err := common.NewPostgresConnection(appConfig.Dbconfig.ConnectionStr)
	if err != nil {
		apilogger.Fatal(
			"error connecting the database",
			zap.Error(err),
			zap.String("connection_string", appConfig.Dbconfig.ConnectionStr),
		)
	}

	productRepository, err := product.NewProductRepositoryWithUser(sqlConnection)
	if err != nil {
		apilogger.Fatal(
			"Something went wrong when initializing the repository",
			zap.Error(err),
		)
	}

	productService := product.NewProductService(productRepository)

	productRouter := product.NewProductRouter(productService)

	app := application.New(productRouter)

	server := &http.Server{
		Addr:    appConfig.Port,
		Handler: app.Mount(),
	}

	apilogger.Debug(
		"configuration loaded",
		zap.String("port", appConfig.Port),
		zap.String("db", appConfig.Dbconfig.ConnectionStr),
		zap.String("rabbitmq_connection", appConfig.Rabbitmqconfig.ConnectionStr),
	)

	apilogger.Info(
		"server running",
		zap.String("port", appConfig.Port),
	)

	eventPublisher, err := events.NewRabbitMqConnection(
		appConfig.Rabbitmqconfig.ConnectionStr,
		appConfig.Rabbitmqconfig.ExchangeName,
	)

	if err != nil {
		apilogger.Fatal(
			"error initializing event publisher",
			zap.String("error", err.Error()),
		)
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

	apilogger.Fatal("server error", zap.String("error", err.Error()))

}
