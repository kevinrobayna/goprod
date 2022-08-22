package web

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kevinrobayna/goprod/internal"
	"github.com/stretchr/testify/assert"
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
	t.Parallel()

	t.Run("Hello", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		data, err := internal.TestWithPostgres(ctx)
		assert.NoError(t, err)
		defer func(postgresC testcontainers.Container, ctx context.Context) {
			err := data.Container.Terminate(ctx)
			if err != nil {
				t.Error(err)
			}
		}(data.Container, ctx)

		var port Port
		app := fxtest.New(t, internal.TestModule, Module, fx.Replace(data.Config), fx.Populate(&port))
		defer app.RequireStart().RequireStop()

		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/", port))
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
		data, err := internal.TestWithPostgres(ctx)
		assert.NoError(t, err)
		defer func(postgresC testcontainers.Container, ctx context.Context) {
			err := data.Container.Terminate(ctx)
			if err != nil {
				t.Error(err)
			}
		}(data.Container, ctx)

		var port Port
		app := fxtest.New(t, internal.TestModule, Module, fx.Replace(data.Config), fx.Populate(&port))
		defer app.RequireStart().RequireStop()

		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/ping", port))
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		var p []internal.Product
		err = json.Unmarshal(body, &p)
		assert.NoError(t, err)

		assert.Equal(t, 0, len(p))
	})
}
