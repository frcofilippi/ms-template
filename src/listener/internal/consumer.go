package internal

import (
	"context"
	"encoding/json"
	"frcofilippi/pedimeapp/shared/config"
	"frcofilippi/pedimeapp/shared/events"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn           *amqp.Connection
	queueName      string
	messageHandler MessageHanlder
}

func (c *Consumer) ListenForMessages(ctx context.Context) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	messageChan, err := ch.ConsumeWithContext(ctx, c.queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			log.Default().Printf("Context canceled.")
			return ctx.Err()
		default:
			for msg := range messageChan {
				var parsedMessage events.OutboxMessage
				body := msg.Body
				err := json.Unmarshal(body, &parsedMessage)
				if err != nil {
					log.Default().Printf("[Consumer] - Error parsing message \n")
					continue
				}
				log.Default().Printf(
					"[Consumer] - Message received. Id: %d Type: %s \n",
					parsedMessage.Id,
					parsedMessage.EventType,
				)
				errCh := make(chan error)
				go c.messageHandler.Handle(parsedMessage, errCh)
				err = <-errCh
				close(errCh)
				if err != nil {
					log.Default().Printf("[Consumer] - Error consumming the message. Erro: %s", err.Error())
					msg.Nack(false, false)
					continue
				}
				msg.Ack(false)
			}
		}
	}
}

func NewConsumer(connection *amqp.Connection, appConfig config.RabbitMQConfiguration, handler MessageHanlder) (*Consumer, error) {

	consummer := &Consumer{
		conn:           connection,
		messageHandler: handler,
	}

	err := consummer.setup(appConfig.ExchangeName, appConfig.DLExchange, appConfig.DLQueue)
	if err != nil {
		return nil, err
	}

	return consummer, nil
}

func (c *Consumer) setup(exchangeName, dlExchangeName, dlQueueName string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(
		dlExchangeName, // DLX name
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err

	}

	err = declareDeadLetterQueue(ch, dlQueueName, dlExchangeName)
	if err != nil {
		return err
	}

	q, err := declareRandomQueue(ch, exchangeName, dlExchangeName)
	if err != nil {
		return err
	}
	c.queueName = q.Name

	return nil
}

func declareDeadLetterQueue(ch *amqp.Channel, dlQueueName string, dlExchangeName string) error {
	_, err := ch.QueueDeclare(
		dlQueueName,
		true, // durable
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	err = ch.QueueBind(
		dlQueueName,
		"",
		dlExchangeName, // bind to DLX
		false,
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}

func declareRandomQueue(ch *amqp.Channel, exchangeBind, deadletterExchangeBind string) (amqp.Queue, error) {
	q, err := ch.QueueDeclare("", false, false, true, false, amqp.Table{
		"x-dead-letter-exchange": deadletterExchangeBind,
	})
	if err != nil {
		return amqp.Queue{}, err
	}

	err = ch.QueueBind(q.Name, "", exchangeBind, false, nil)
	if err != nil {
		return amqp.Queue{}, err
	}

	return q, nil
}
