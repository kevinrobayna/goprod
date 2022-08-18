package main

import (
	"context"
	"encoding/json"
	"github.com/kevinrobayna/goprod/internal"
	"github.com/kevinrobayna/goprod/internal/di"
	"github.com/kevinrobayna/goprod/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"io"
	"net/http"
	"testing"
)

type Message struct {
	Message string `json:"message"`
}

func TestRoutes(t *testing.T) {
	t.Run("Hello", func(t *testing.T) {
		ctx := context.Background()

		pg, err := di.TestWithPostgres(ctx)
		assert.NoError(t, err)
		defer func(postgresC testcontainers.Container, ctx context.Context) {
			err := pg.Container.Terminate(ctx)
			if err != nil {
				t.Error(err)
			}
		}(pg.Container, ctx)

		app := fxtest.New(t, opts(), fx.Replace(pg.Config))
		defer app.RequireStart().RequireStop()

		resp, err := http.Get("http://localhost:8080/")
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		var message = Message{}
		err = json.Unmarshal(body, &message)
		assert.NoError(t, err)

		assert.Equal(t, "hello", message.Message)
	})

	t.Run("Ping", func(t *testing.T) {
		ctx := context.Background()
		data, err := di.TestWithPostgres(ctx)
		assert.NoError(t, err)
		defer func(postgresC testcontainers.Container, ctx context.Context) {
			err := data.Container.Terminate(ctx)
			if err != nil {
				t.Error(err)
			}
		}(data.Container, ctx)

		app := fxtest.New(t, opts(), fx.Replace(data.Config))
		defer app.RequireStart().RequireStop()

		resp, err := http.Get("http://localhost:8080/ping")
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		var product = models.Product{}
		err = json.Unmarshal(body, &product)
		assert.NoError(t, err)
		var p []internal.Product
		err = json.Unmarshal(body, &p)
		if err != nil {
			panic(err.Error())
		}

		assert.NotEmpty(t, product)
		assert.Equal(t, "D42", product.Code)
		assert.Equal(t, uint(100), product.Price)
		assert.Equal(t, 0, len(p))
	})
}
