package internal

import (
	"context"
	"encoding/json"
)

type QueueMessage struct {
	MessageId   int64           `json:"message_id"`
	MessageType string          `json:"message_type"`
	Data        json.RawMessage `json:"data"`
}

type Listener interface {
	ListenForMessages(ctx context.Context) error
}

type MessageHanlder interface {
	Handle(QueueMessage, chan<- error)
}
