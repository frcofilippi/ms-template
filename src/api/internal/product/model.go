package product

import (
	"errors"
	"frcofilippi/pedimeapp/shared/events"
	product_events "frcofilippi/pedimeapp/shared/events/product"
	"strings"
	"time"
)

const (
	validationError = "name or cost were not provided"
)

type Product struct {
	Id     int64
	UserId string
	Name   string
	Cost   float64
	events []events.DomainEvent
}

//TODO: decide whether to create a productId or to use the sequence

func (p *Product) RaiseEvent(event events.DomainEvent) {
	p.events = append(p.events, event)
}

func (p *Product) PendingEvents() []events.DomainEvent {
	for _, evt := range p.events {
		prodCreatedEvt := evt.(*product_events.ProductCreatedEvent)
		prodCreatedEvt.ProductId = p.Id
	}
	return p.events
}

func (p *Product) CleanEvents() {
	p.events = make([]events.DomainEvent, 0)
}

func NewProduct(id int64, userId string, name string, cost float64) (*Product, error) {
	if name == "" || cost == 0 {
		return nil, errors.New(validationError)
	}
	if userId == "" {
		return nil, errors.New("user must be set before creating a product")
	}
	product := &Product{
		Id:     id,
		UserId: userId,
		Name:   strings.ToUpper(name),
		Cost:   cost,
	}

	productCreatedEvent := &product_events.ProductCreatedEvent{
		ProductId: product.Id,
		Cost:      product.Cost,
		Time:      time.Now(),
	}

	product.RaiseEvent(productCreatedEvent)

	return product, nil
}
