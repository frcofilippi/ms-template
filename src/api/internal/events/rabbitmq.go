package events

import (
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Exchange   string
}

func NewRabbitMqConnection(amqpURL, exchangeName string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Declare exchange
	err = ch.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Guarantee messages are delivered one at a time per consumer
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	return &RabbitMQ{
		Connection: conn,
		Channel:    ch,
		Exchange:   exchangeName,
	}, nil
}

// Publish sends a message to the exchange.
func (r *RabbitMQ) Publish(routingKey string, body []byte) error {
	return r.Channel.Publish(
		r.Exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// Close closes the channel and connection.
func (r *RabbitMQ) Close() error {
	if r.Channel != nil {
		_ = r.Channel.Close()
	}
	if r.Connection != nil {
		_ = r.Connection.Close()
	}
	return nil
}
