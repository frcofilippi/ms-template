package main

type RabbitMQConfiguration struct {
	connectionStr string
	exchangeName  string
	queue         string
}

func NewRabbitMQConfig() *RabbitMQConfiguration {
	return &RabbitMQConfiguration{
		connectionStr: readEnvValueAsString("RMQ_CONNECTION_STR", ""),
		exchangeName:  readEnvValueAsString("RMQ_EXCHANGE", "pedimeapp_exchange"),
		queue:         readEnvValueAsString("RMQ_QUEUE", "pedimeapp_queue"),
	}
}
