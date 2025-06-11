package product

import (
	"context"
	"database/sql"
)

type CreateProductRequest struct {
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}

type WithUserProductRepository struct {
	innerRepo *PgProductRepository
	db        *sql.DB
}

type DBExecutor interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type ProductService interface {
	GetProductById(GetProductByIdCommand) (*Product, error)
	CreateNewProduct(CreateNewProductCommand) (int64, error)
}

type GetProductByIdCommand struct {
	ProductId int64
	UserId    string
}

type CreateNewProductCommand struct {
	Name   string
	Cost   float64
	UserId string
}

type ProductRepository interface {
	GetById(context.Context, int64, string) (*Product, error)
	Create(context.Context, *Product) (int64, error)
}
