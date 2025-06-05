package application

import (
	"context"
	"encoding/json"
	"frcofilippi/pedimeapp/internal/events"

	"github.com/streadway/amqp"
)

type MessageListener struct {
	consummer  events.EventConsumer
	dispatcher events.MessageDispatcher
}

func (ml *MessageListener) Start() error {
	return ml.consummer.Consume(func(d amqp.Delivery) {
		var message OutboxMessage
		err := json.Unmarshal(d.Body, &message)
		if err != nil {
			d.Acknowledger.Nack(d.DeliveryTag, false, true)
			return
			// errChan <- err
		}
		err = ml.dispatcher.Dispatch(context.Background(), message.EventType, message.Payload)
		if err != nil {
			d.Acknowledger.Nack(d.DeliveryTag, false, true)
			return
			// errChan <- err
		}
		err = d.Acknowledger.Ack(d.DeliveryTag, false)
		if err != nil {
			// errChan <- err
			d.Acknowledger.Nack(d.DeliveryTag, false, true)
			return
		}
	})
}

func NewMessageListener(consummer events.EventConsumer, dispatcher events.MessageDispatcher) *MessageListener {
	return &MessageListener{
		consummer:  consummer,
		dispatcher: dispatcher,
	}
}
