package main

import (
	"frcofilippi/pedimeapp/internal/application"
	"frcofilippi/pedimeapp/internal/application/services"
	"frcofilippi/pedimeapp/internal/infra/database"
	"frcofilippi/pedimeapp/internal/infra/handlers"
	"frcofilippi/pedimeapp/internal/infra/repositories"
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
	sqlConnection, err := database.NewPostgresConnection(appConfig.dbconfig.connectionStr)
	if err != nil {
		log.Fatalf("error connecting the database. Error %v", err)
	}

	// productRepository, err := repositories.NewProductRepository(sqlConnection)

	productRepository, err := repositories.NewProductRepositoryWithCustomer(sqlConnection)
	if err != nil {
		log.Fatalf("Something went wrong when initializing the repository. Error: %v", err)
	}

	productService := services.NewProductService(productRepository)

	productRouter := handlers.NewProductRouter(productService)

	app := application.New(productRouter)

	server := &http.Server{
		Addr:    appConfig.port,
		Handler: app.Mount(),
	}

	log.Printf("Configuration Loaded. Port: %s - Db: %s \n", appConfig.port, appConfig.dbconfig.connectionStr)
	log.Printf("Server running on port: %s \n", appConfig.port)

	err = server.ListenAndServe()

	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

}
