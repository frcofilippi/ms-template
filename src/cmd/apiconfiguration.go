package main

import (
	"os"
	"strconv"
)

type ApiConfiguration struct {
	port           string
	dbconfig       *DatabaseConfiguration
	rabbitmqconfig *RabbitMQConfiguration
}

func readEnvValueAsString(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func readEnvValueAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	oValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return oValue
}

func NewApiConfiguration() *ApiConfiguration {
	return &ApiConfiguration{
		port:           readEnvValueAsString("APP_PORT", ":2020"),
		dbconfig:       NewDatabaseConfig(),
		rabbitmqconfig: NewRabbitMQConfig(),
	}
}
