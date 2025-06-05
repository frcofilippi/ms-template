package application

import (
	"context"
	"encoding/json"
	"log"
)

const eventTypeRegistered = "ProductCreated"

type MessageDispatcher struct {
}

func (md *MessageDispatcher) Dispatch(ctx context.Context, enventType string, payload json.RawMessage) error {
	if enventType == eventTypeRegistered {
		log.Default().Printf("[DISPATCHER] - Handling message %v \n", string(payload))
	}
	return nil
}
