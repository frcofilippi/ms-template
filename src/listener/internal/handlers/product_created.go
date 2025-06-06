package handlers

import (
	"encoding/json"
	events "frcofilippi/pedimeapp/shared/events"
	product_events "frcofilippi/pedimeapp/shared/events/product"
	"log"
)

type ProductCreatedHandler struct{}

func (pc *ProductCreatedHandler) Handle(message events.OutboxMessage, errChan chan<- error) {
	var pcEvent product_events.ProductCreatedEvent
	err := json.Unmarshal(message.Payload, &pcEvent)
	if err != nil {
		errChan <- err
		return
	}
	log.Default().Printf("[ProductCreatedHandler] - Event processed - %v \n", pcEvent)
	errChan <- nil
}
