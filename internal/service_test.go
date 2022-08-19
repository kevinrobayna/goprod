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

		args := product{Code: "code", Description: "description", Price: 1}
		repo.On("SaveProduct", args).Return(randomProduct(), nil)

		service := Service{repo}

		err := service.SaveProduct(Product{Code: "code", Description: "description", Price: 1})
		assert.NoError(t, err)

		repo.AssertExpectations(t)
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
