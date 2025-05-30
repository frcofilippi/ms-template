package events

import (
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
