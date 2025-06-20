package product

import (
	"time"
)

type ProductCreatedEvent struct {
	ProductId int64     `json:"product_id"`
	Cost      float64   `json:"cost"`
	Time      time.Time `json:"created_at"`
}

const ProductCreatedEventName = "ProductCreated"

func (pc *ProductCreatedEvent) EventType() string {
	return ProductCreatedEventName
}

func (pc *ProductCreatedEvent) OccurredAt() time.Time {
	return pc.Time
}
