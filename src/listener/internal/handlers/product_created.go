package handlers

import (
	"context"
	"encoding/json"
	events "frcofilippi/pedimeapp/shared/events"
	product_events "frcofilippi/pedimeapp/shared/events/product"
	"frcofilippi/pedimeapp/shared/logger"

	"go.uber.org/zap"
)

type ProductCreatedHandler struct {
	HandledEvent string
}

func (pc *ProductCreatedHandler) HandleContext(ctx context.Context, message events.OutboxMessage) error {
	var pcEvent product_events.ProductCreatedEvent
	err := json.Unmarshal(message.Payload, &pcEvent)
	if err != nil {
		return err
	}
	logger.Info("Event processed", zap.String("eventName", pc.HandledEvent), zap.Any("event_data", pcEvent))
	return nil
}

func NewProductCreatedEventHandler(eventHandled string) *ProductCreatedHandler {
	return &ProductCreatedHandler{
		HandledEvent: eventHandled,
	}
}
