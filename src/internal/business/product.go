package business

import (
	"errors"
	"strings"
)

const (
	validationError = "name or cost were not provided"
)

type Product struct {
	Id         int64
	CustomerId int64
	Name       string
	Cost       float64
}

func NewProduct(id, customer int64, name string, cost float64) (*Product, error) {
	if name == "" || cost == 0 {
		return nil, errors.New(validationError)
	}
	if customer == 0 {
		return nil, errors.New("customer must be set before creating a product")
	}
	return &Product{
		Id:         id,
		CustomerId: customer,
		Name:       strings.ToUpper(name),
		Cost:       cost,
	}, nil
}
