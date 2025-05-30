package events

import "time"

type ProductCreatedEvent struct {
	ProductId int64     `json:"product_id"`
	Cost      float64   `json:"cost"`
	Time      time.Time `json:"created_at"`
}

const EventName = "ProductCreated"

func (pc *ProductCreatedEvent) EventType() string {
	return EventName
}

func (pc *ProductCreatedEvent) OccurredAt() time.Time {
	return pc.Time
}
