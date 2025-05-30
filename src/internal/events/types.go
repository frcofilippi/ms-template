package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/streadway/amqp"
)

type EventPublisher interface {
	Publish(routingKey string, body []byte) error
	Close() error
}

type EventConsumer interface {
	Consume(handler func(amqp.Delivery)) error
	Close() error
}

type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
}

type MessageDispatcher interface {
	Dispatch(ctx context.Context, enventType string, payload json.RawMessage) error
}
