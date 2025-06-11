package product

import (
	"context"
	"errors"
)

type ApplicationProductService struct {
	repository ProductRepository
}

func (ps *ApplicationProductService) GetProductById(cmd GetProductByIdCommand) (*Product, error) {
	if cmd.ProductId == 0 {
		return nil, errors.New("id must be provided to find the product")
	}

	prod, err := ps.repository.GetById(context.Background(), cmd.ProductId, cmd.UserId)
	if err != nil {
		return nil, err
	}
	return prod, nil
}

func (ps *ApplicationProductService) CreateNewProduct(cmd CreateNewProductCommand) (int64, error) {
	if cmd.UserId == "" {
		return 0, errors.New("customer id must be set to create a product")
	}

	prod, err := NewProduct(0, cmd.UserId, cmd.Name, cmd.Cost)

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
