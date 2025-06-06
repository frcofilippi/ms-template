package internal

import (
	"fmt"
	"frcofilippi/pedimeapp/shared/events"
	"log"
)

type CustomMessageHandler struct{}

func (handler *CustomMessageHandler) Handle(message events.OutboxMessage, errCh chan<- error) {
	log.Default().Printf("[CustomHandler] - Handling message. Message id: %d", message.Id)
	if message.EventType == "Unknown" {
		errCh <- fmt.Errorf("unknown message type - message_id: %d", message.Id)
		return
	}
	errCh <- nil
}
