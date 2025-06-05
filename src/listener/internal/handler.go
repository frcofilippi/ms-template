package internal

import (
	"fmt"
	"log"
)

type CustomMessageHandler struct{}

func (handler *CustomMessageHandler) Handle(message QueueMessage, errCh chan<- error) {
	log.Default().Printf("[CustomHandler] - Handling message. Message id: %d", message.MessageId)
	if message.MessageType == "Unknown" {
		errCh <- fmt.Errorf("unknown message type - message_id: %d", message.MessageId)
		return
	}
	errCh <- nil
}
