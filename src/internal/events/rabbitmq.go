package events

import (
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
	Exchange   string
}

func NewRabbitMqConnection(amqpURL, exchangeName, mainqName, dlxName, dlqName string) (*RabbitMQ, error) {
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

	//declare deadlatter exchange
	err = ch.ExchangeDeclare(
		dlqName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	// Declare queue
	q, err := ch.QueueDeclare(
		mainqName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		q.Name,       // queue name
		"",           // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind queue: %w", err)
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
		Queue:      q,
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

// Consume starts consuming messages and calls handler for each delivery.
func (r *RabbitMQ) Consume(handler func(amqp.Delivery)) error {
	msgs, err := r.Channel.Consume(
		r.Queue.Name,
		"",
		false, // auto-ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for d := range msgs {
			handler(d)
		}
	}()

	return nil
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
