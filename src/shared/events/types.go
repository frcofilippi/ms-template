package events

import (
	"encoding/json"
	"time"

	"github.com/streadway/amqp"
)

// OutboxMessage represents the structure of messages passed through the queue
type OutboxMessage struct {
	Id        int64           `json:"id"`
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

// EventPublisher defines the contract for publishing events
type EventPublisher interface {
	Publish(routingKey string, body []byte) error
	Close() error
}

// EventConsumer defines the contract for consuming events
type EventConsumer interface {
	Consume(handler func(amqp.Delivery)) error
	Close() error
}

// MessageHandler defines how messages should be handled
type MessageHandler interface {
	Handle(OutboxMessage, chan<- error)
}

// Event represents any event in the system
type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
}
