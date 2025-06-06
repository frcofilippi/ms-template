package internal

import (
	"context"
	"frcofilippi/pedimeapp/shared/events"
)

// type OutboxMessage struct {
// 	Id          int64           `json:"id"`
// 	MessageType string          `json:"event_type"`
// 	Payload     json.RawMessage `json:"payload"`
// }

type Listener interface {
	ListenForMessages(ctx context.Context) error
}

type MessageHanlder interface {
	Handle(events.OutboxMessage, chan<- error)
}
