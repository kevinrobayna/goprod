package internal

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"math/rand"
	"testing"
	"time"
)

type MyMockedObject struct {
	mock.Mock
}

func (m *MyMockedObject) SaveProduct(p product) (product, error) {
	args := m.Called(p)
	return args.Get(0).(product), args.Error(1)
}

func (m *MyMockedObject) UpdateProduct(p product) (product, error) {
	args := m.Called(p)
	return args.Get(0).(product), args.Error(1)
}

func (m *MyMockedObject) GetProducts() ([]product, error) {
	args := m.Called()
	return args.Get(0).([]product), args.Error(1)
}

func (m *MyMockedObject) GetProduct(id uint) (product, error) {
	args := m.Called(id)
	return args.Get(0).(product), args.Error(1)
}

func TestService(t *testing.T) {
	t.Parallel()

	t.Run("Save New Product", func(t *testing.T) {
		t.Parallel()

		repo := new(MyMockedObject)
		service := Service{repo}

		repo.On("SaveProduct", mock.Anything).Return(randomProduct(), nil)

		err := service.SaveProduct(Product{Code: "code", Description: "description", Price: 1})
		assert.NoError(t, err)

		repo.AssertCalled(t, "SaveProduct", mock.Anything)
		repo.AssertNotCalled(t, "UpdateProduct")
	})

	t.Run("Update existing product", func(t *testing.T) {
		t.Parallel()

		repo := new(MyMockedObject)
		service := Service{repo}

		p := randomProduct()
		repo.On("GetProduct", p.ID).Return(p, nil)

		repo.On("UpdateProduct", mock.Anything).Return(p, nil)

		product := Product{Id: p.ID, Code: "code", Description: "description", Price: 1}
		err := service.SaveProduct(product)
		assert.NoError(t, err)

		repo.AssertCalled(t, "UpdateProduct", mock.Anything)
		repo.AssertNotCalled(t, "SaveProduct")
	})

	t.Run("Get Existing Repo", func(t *testing.T) {
		t.Parallel()

		repo := new(MyMockedObject)
		service := Service{repo}

		p := randomProduct()
		repo.On("GetProduct", p.ID).Return(p, nil)

		product, err := service.GetProduct(p.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, product)

		repo.AssertCalled(t, "GetProduct", p.ID)

		assert.Equal(t, p.Code, product.Code)
		assert.Equal(t, p.Description, product.Description)
		assert.Equal(t, p.Price, product.Price)
	})

	t.Run("Get not existing product", func(t *testing.T) {
		t.Parallel()

		repo := new(MyMockedObject)
		service := Service{repo}

		repo.On("GetProduct", uint(1)).Return(product{}, nil)

		product, err := service.GetProduct(uint(1))
		assert.NoError(t, err)
		assert.Empty(t, product)

		repo.AssertCalled(t, "GetProduct", uint(1))
	})

	t.Run("Get all when nothing exist", func(t *testing.T) {
		t.Parallel()

		repo := new(MyMockedObject)
		service := Service{repo}

		repo.On("GetProducts").Return([]product{}, nil)

		products, err := service.GetProducts()
		assert.NoError(t, err)
		assert.Empty(t, products)

		repo.AssertCalled(t, "GetProducts")
	})

	t.Run("Get all products", func(t *testing.T) {
		t.Parallel()

		repo := new(MyMockedObject)
		service := Service{repo}

		var p []product
		p = append(p, randomProduct())
		p = append(p, randomProduct())
		p = append(p, randomProduct())
		repo.On("GetProducts").Return(p, nil)

		products, err := service.GetProducts()
		assert.NoError(t, err)
		assert.NotEmpty(t, products)
		assert.Len(t, products, 3)

		repo.AssertCalled(t, "GetProducts")
	})
}

func randomProduct() product {
	return product{
		Model:       gorm.Model{ID: uint(rand.Int()), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Code:        "code",
		Description: "description",
		Price:       1,
	}
}
