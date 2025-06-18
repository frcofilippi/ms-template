package product

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetById(ctx context.Context, productId int64, userId string) (*Product, error) {
	args := m.Called(ctx, productId, userId)
	if prod, ok := args.Get(0).(*Product); ok {
		return prod, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepository) Create(ctx context.Context, product *Product) (int64, error) {
	args := m.Called(ctx, product)
	return args.Get(0).(int64), args.Error(1)
}

func TestGetProductById_CustomerIdSet(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)
	cmd := GetProductByIdCommand{ProductId: 1, UserId: "user1"}
	expectedProduct := &Product{Id: 1, UserId: "user1", Name: "Test", Cost: 10}

	mockRepo.On("GetById", mock.Anything, int64(1), "user1").Return(expectedProduct, nil)

	prod, err := service.GetProductById(cmd)
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, prod)
	mockRepo.AssertExpectations(t)
}

func TestGetProductById_CustomerIdNotSet(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)
	cmd := GetProductByIdCommand{ProductId: 0, UserId: ""}

	prod, err := service.GetProductById(cmd)
	assert.Nil(t, prod)
	assert.EqualError(t, err, "id must be provided to find the product")
}

// func TestCreateNewProduct_CustomerIdSet(t *testing.T) {
// 	mockRepo := new(MockProductRepository)
// 	service := NewProductService(mockRepo)
// 	cmd := CreateNewProductCommand{UserId: "user1", Name: "Test", Cost: 10}
// 	// expectedProduct := &Product{Id: 0, UserId: "user1", Name: "Test", Cost: 10}

// 	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *Product) bool {
// 		return p.UserId == "user1" && p.Name == "Test" && p.Cost == 10
// 	})).Return(int64(123), nil)

// 	id, err := service.CreateNewProduct(cmd)
// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(123), id)
// 	mockRepo.AssertExpectations(t)
// }

func TestCreateNewProduct_CustomerIdNotSet(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)
	cmd := CreateNewProductCommand{UserId: "", Name: "Test", Cost: 10}

	id, err := service.CreateNewProduct(cmd)
	assert.Equal(t, int64(0), id)
	assert.EqualError(t, err, "customer id must be set to create a product")
}
