package product

import (
	"errors"
	"frcofilippi/pedimeapp/internal/events"
	"strings"
	"time"
)

const (
	validationError = "name or cost were not provided"
)

type Product struct {
	Id         int64
	CustomerId int64
	Name       string
	Cost       float64
	events     []events.DomainEvent
}

//TODO: decide whether to create a productId or to use the sequence

func (p *Product) RaiseEvent(event events.DomainEvent) {
	p.events = append(p.events, event)
}

func (p *Product) PendingEvents() []events.DomainEvent {
	for _, evt := range p.events {
		prodCreatedEvt := evt.(*events.ProductCreatedEvent)
		prodCreatedEvt.ProductId = p.Id
	}
	return p.events
}

func (p *Product) CleanEvents() {
	p.events = make([]events.DomainEvent, 0)
}

func NewProduct(id, customer int64, name string, cost float64) (*Product, error) {
	if name == "" || cost == 0 {
		return nil, errors.New(validationError)
	}
	if customer == 0 {
		return nil, errors.New("customer must be set before creating a product")
	}
	product := &Product{
		Id:         id,
		CustomerId: customer,
		Name:       strings.ToUpper(name),
		Cost:       cost,
	}

	productCreatedEvent := &events.ProductCreatedEvent{
		ProductId: product.Id,
		Cost:      product.Cost,
		Time:      time.Now(),
	}

	product.RaiseEvent(productCreatedEvent)

	return product, nil
}
