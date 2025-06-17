package events

import (
	"context"
	"encoding/json"
	"time"
)

type OutboxMessage struct {
	Id        int64           `json:"id"`
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

type RabbitMQConfig struct {
	URL        string
	Exchange   string
	DLExchange string
	DLQueue    string
}

type HandlerFunc func(ctx context.Context, msg OutboxMessage) error

type Dispatcher map[string]HandlerFunc

type EventPublisher interface {
	Publish(ctx context.Context, routingKey string, body []byte) error
	Close() error
}

type EventConsumer interface {
	Listen(ctx context.Context, dispatcher Dispatcher) error
	Close() error
}

type MessageHandler interface {
	Handle(OutboxMessage, chan<- error)
}

type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
}
