package config

import (
	"os"
)

type RabbitMQConfiguration struct {
	ConnectionStr string
	ExchangeName  string
	DLExchange    string
	DLQueue       string
}

type ApiConfiguration struct {
	Port           string
	Dbconfig       *DatabaseConfiguration
	Rabbitmqconfig *RabbitMQConfiguration
}

func NewApiConfiguration() *ApiConfiguration {
	// err := godotenv.Load("../../.env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	return &ApiConfiguration{
		Port:           readEnvValueAsString("APP_PORT", ":2020"),
		Dbconfig:       NewDatabaseConfig(),
		Rabbitmqconfig: NewRabbitMQConfig(),
	}
}

func readEnvValueAsString(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func NewRabbitMQConfig() *RabbitMQConfiguration {
	return &RabbitMQConfiguration{
		ConnectionStr: readEnvValueAsString("RMQ_CONNECTION_STR", ""),
		ExchangeName:  readEnvValueAsString("RMQ_EXCHANGE", "pedimeapp_exchange"),
		DLQueue:       readEnvValueAsString("RMQ_DL_QUEUE", "pedimeapp_dl_queue"),
		DLExchange:    readEnvValueAsString("RMQ_DL_EXCHANGE", "pedimeapp_dl_exchange"),
	}
}

type DatabaseConfiguration struct {
	ConnectionStr string
}

func NewDatabaseConfig() *DatabaseConfiguration {
	return &DatabaseConfiguration{
		ConnectionStr: readEnvValueAsString("DB_CONNECTION_STR", "postgres://pedimeapp:mysecretpwd@localhost/pedimedb?sslmode=disable"),
	}
}
