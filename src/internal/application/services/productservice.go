package services

import (
	"context"
	"errors"
	"frcofilippi/pedimeapp/internal/business"
)

type ProductService interface {
	GetProductById(GetProductByIdCommand) (*business.Product, error)
	CreateNewProduct(CreateNewProductCommand) (int64, error)
}

type GetProductByIdCommand struct {
	ProductId  int64
	CustomerId int64
}

type CreateNewProductCommand struct {
	Name       string
	Cost       float64
	CustomerId int64
}

type ProductRepository interface {
	GetById(context.Context, int64, int64) (*business.Product, error)
	Create(context.Context, *business.Product) (int64, error)
}

type ApplicationProductService struct {
	repository ProductRepository
}

func (ps *ApplicationProductService) GetProductById(cmd GetProductByIdCommand) (*business.Product, error) {
	if cmd.ProductId == 0 {
		return nil, errors.New("id must be provided to find the product")
	}

	prod, err := ps.repository.GetById(context.Background(), cmd.ProductId, cmd.CustomerId)
	if err != nil {
		return nil, err
	}
	return prod, nil
}

func (ps *ApplicationProductService) CreateNewProduct(cmd CreateNewProductCommand) (int64, error) {
	if cmd.CustomerId == 0 {
		return 0, errors.New("customer id must be set to create a product")
	}

	prod, err := business.NewProduct(0, cmd.CustomerId, cmd.Name, cmd.Cost)

	if err != nil {
		return 0, err
	}

	createdProductId, err := ps.repository.Create(context.Background(), prod)
	if err != nil {
		return 0, err
	}

	return createdProductId, nil
}

func NewProductService(repo ProductRepository) *ApplicationProductService {
	return &ApplicationProductService{
		repository: repo,
	}
}
