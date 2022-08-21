package internal

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"math/rand"
	"testing"
)

func TestRepository(t *testing.T) {
	t.Parallel()

	t.Run("SaveProduct", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		data, err := TestWithPostgres(ctx)
		assert.NoError(t, err)
		defer func(postgresC testcontainers.Container, ctx context.Context) {
			err := data.Container.Terminate(ctx)
			if err != nil {
				t.Error(err)
			}
		}(data.Container, ctx)

		var repository iRepository
		app := fxtest.New(t, TestModule, Module, fx.NopLogger, fx.Replace(data.Config), fx.Populate(&repository))
		defer app.RequireStart().RequireStop()

		product := createRandomProduct()
		savedProduct, err := repository.SaveProduct(product)
		assert.NoError(t, err)

		assert.Empty(t, product.ID)
		assert.NotEmpty(t, savedProduct.ID)

		foundProduct, err := repository.GetProduct(savedProduct.ID)
		assert.NoError(t, err)

		assert.NotEmpty(t, foundProduct)
	})

	t.Run("SaveProductTwiceFails", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		data, err := TestWithPostgres(ctx)
		assert.NoError(t, err)
		defer func(postgresC testcontainers.Container, ctx context.Context) {
			err := data.Container.Terminate(ctx)
			if err != nil {
				t.Error(err)
			}
		}(data.Container, ctx)

		var repository iRepository
		app := fxtest.New(t, TestModule, Module, fx.NopLogger, fx.Replace(data.Config), fx.Populate(&repository))
		defer app.RequireStart().RequireStop()

		product, err := repository.SaveProduct(createRandomProduct())
		assert.NoError(t, err)

		_, err = repository.SaveProduct(product)
		assert.Error(t, err)
	})

	t.Run("UpdateProduct", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		data, err := TestWithPostgres(ctx)
		assert.NoError(t, err)
		defer func(postgresC testcontainers.Container, ctx context.Context) {
			err := data.Container.Terminate(ctx)
			if err != nil {
				t.Error(err)
			}
		}(data.Container, ctx)

		var repository iRepository
		app := fxtest.New(t, TestModule, Module, fx.NopLogger, fx.Replace(data.Config), fx.Populate(&repository))
		defer app.RequireStart().RequireStop()

		product, err := repository.SaveProduct(createRandomProduct())
		assert.NoError(t, err)

		product.Description = "Updated"
		updatedProduct, err := repository.UpdateProduct(product)
		assert.NoError(t, err)

		assert.Equal(t, "Updated", updatedProduct.Description)
	})

	t.Run("GetProduct", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		data, err := TestWithPostgres(ctx)
		assert.NoError(t, err)
		defer func(postgresC testcontainers.Container, ctx context.Context) {
			err := data.Container.Terminate(ctx)
			if err != nil {
				t.Error(err)
			}
		}(data.Container, ctx)

		var repository iRepository
		app := fxtest.New(t, TestModule, Module, fx.NopLogger, fx.Replace(data.Config), fx.Populate(&repository))
		defer app.RequireStart().RequireStop()

		product, err := repository.SaveProduct(createRandomProduct())
		assert.NoError(t, err, "Unable to save product in db")

		savedProduct, err := repository.GetProduct(product.ID)
		assert.NoError(t, err, "Unable to find product in db")

		assert.NotEmpty(t, savedProduct.Price)
		assert.Equal(t, product.Price, savedProduct.Price)
		assert.Equal(t, product.Description, savedProduct.Description)
		assert.Equal(t, product.Code, savedProduct.Code)
	})

	t.Run("GetProducts", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		data, err := TestWithPostgres(ctx)
		assert.NoError(t, err)
		defer func(postgresC testcontainers.Container, ctx context.Context) {
			err := data.Container.Terminate(ctx)
			if err != nil {
				t.Error(err)
			}
		}(data.Container, ctx)

		var repository iRepository
		app := fxtest.New(t, TestModule, Module, fx.NopLogger, fx.Replace(data.Config), fx.Populate(&repository))
		defer app.RequireStart().RequireStop()

		_, err = repository.SaveProduct(createRandomProduct())
		assert.NoError(t, err, "Unable to save product in db")
		_, err = repository.SaveProduct(createRandomProduct())
		assert.NoError(t, err, "Unable to save product in db")
		_, err = repository.SaveProduct(createRandomProduct())
		assert.NoError(t, err, "Unable to save product in db")

		savedProducts, err := repository.GetProducts()
		assert.NoError(t, err, "Unable to find product in db")

		assert.NotEmpty(t, savedProducts)
		assert.Equal(t, 3, len(savedProducts))
	})
}

func createRandomProduct() product {
	return product{
		Price:       uint(rand.Int()),
		Code:        "Code",
		Description: "Product",
	}
}
