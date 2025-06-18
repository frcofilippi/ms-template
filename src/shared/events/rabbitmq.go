package events

import (
	"context"
	"encoding/json"
	"fmt"
	"frcofilippi/pedimeapp/shared/logger"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQClient struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
	cfg       RabbitMQConfig
}

func NewRabbitMQClient(cfg RabbitMQConfig) (*RabbitMQClient, error) {
	const (
		maxAttempts    = 5
		initialBackoff = time.Second
		maxBackoff     = 16 * time.Second
	)
	conn, err := retryConnect(cfg.URL, maxAttempts, initialBackoff, maxBackoff)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}
	queueName, err := setupExchangesAndQueues(ch, cfg)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}
	return &RabbitMQClient{conn: conn, channel: ch, queueName: queueName, cfg: cfg}, nil
}

func setupExchangesAndQueues(ch *amqp.Channel, cfg RabbitMQConfig) (string, error) {
	// Declare main exchange
	err := ch.ExchangeDeclare(
		cfg.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to declare exchange: %w", err)
	}
	// Declare DLX
	err = ch.ExchangeDeclare(
		cfg.DLExchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to declare DLX: %w", err)
	}
	// Declare DL queue
	_, err = ch.QueueDeclare(
		cfg.DLQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to declare DL queue: %w", err)
	}

	err = ch.QueueBind(cfg.DLQueue, "", cfg.DLExchange, false, nil)
	if err != nil {
		return "", fmt.Errorf("failed to bind DL queue: %w", err)
	}

	q, err := ch.QueueDeclare("", false, false, true, false, amqp.Table{
		"x-dead-letter-exchange": cfg.DLExchange,
	})
	if err != nil {
		return "", fmt.Errorf("failed to declare main queue: %w", err)
	}
	err = ch.QueueBind(q.Name, "", cfg.Exchange, false, nil)
	if err != nil {
		return "", fmt.Errorf("failed to bind main queue: %w", err)
	}
	return q.Name, nil
}

func (r *RabbitMQClient) Publish(ctx context.Context, routingKey string, body []byte) error {
	return r.channel.PublishWithContext(ctx,
		r.cfg.Exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitMQClient) Listen(ctx context.Context, dispatcher Dispatcher) error {
	msgs, err := r.channel.ConsumeWithContext(ctx, r.queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				return nil
			}
			var outboxMsg OutboxMessage
			err := json.Unmarshal(msg.Body, &outboxMsg)
			if err != nil {
				msg.Nack(false, false)
				continue
			}
			logger.Debug("message received", zap.String("messageType", outboxMsg.EventType))
			handler, found := dispatcher[outboxMsg.EventType]
			if !found {
				msg.Nack(false, false)
				continue
			}
			hErr := handler(ctx, outboxMsg)
			if hErr != nil {
				msg.Nack(false, false)
				continue
			}
			msg.Ack(false)
		}
	}
}

func (r *RabbitMQClient) Close() error {
	if r.channel != nil {
		_ = r.channel.Close()
	}
	if r.conn != nil {
		_ = r.conn.Close()
	}
	return nil
}

func retryConnect(url string, maxAttempts int, initialBackoff, maxBackoff time.Duration) (*amqp.Connection, error) {
	var (
		conn    *amqp.Connection
		err     error
		backoff = initialBackoff
	)
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			return conn, nil
		}
		if attempt < maxAttempts {
			time.Sleep(backoff)
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}
	return nil, err
}
