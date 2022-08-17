package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kevinrobayna/goprod/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"net/http"
	"testing"
)

type Message struct {
	Message string `json:"message"`
}

func SetupWithPostgres() testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image: "postgres:14.5",
		Env: map[string]string{
			"POSTGRES_USER":             "goprod",
			"POSTGRES_PASSWORD":         "",
			"POSTGRES_DB":               "goprod",
			"POSTGRES_HOST_AUTH_METHOD": "trust",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}
}

func TestRoutes(t *testing.T) {
	t.Parallel()

	t.Run("Hello", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		req := SetupWithPostgres()
		postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		assert.NoError(t, err)
		defer postgresC.Terminate(ctx)

		port, err := postgresC.MappedPort(ctx, "5432/tcp")

		dsn := fmt.Sprintf("host=localhost user=goprod password= dbname=goprod port=%s sslmode=disable", port.Port())
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		assert.NoError(t, err)

		app := fxtest.New(t, opts(), fx.Replace(db))
		// Starts the app right away, and defers a stop when the test ends.
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
		t.Parallel()

		ctx := context.Background()
		req := SetupWithPostgres()
		postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		assert.NoError(t, err)
		defer postgresC.Terminate(ctx)

		port, err := postgresC.MappedPort(ctx, "5432/tcp")

		dsn := fmt.Sprintf("host=localhost user=goprod password= dbname=goprod port=%s sslmode=disable", port.Port())
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		assert.NoError(t, err)

		app := fxtest.New(t, opts(), fx.Replace(db))
		// Starts the app right away, and defers a stop when the test ends.
		defer app.RequireStart().RequireStop()

		resp, err := http.Get("http://localhost:8080/ping")
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		var product = models.Product{}
		err = json.Unmarshal(body, &product)
		assert.NoError(t, err)

		assert.NotEmpty(t, product)
		assert.Equal(t, "D42", product.Code)
		assert.Equal(t, uint(100), product.Price)
	})
}
