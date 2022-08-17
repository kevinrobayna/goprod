package main

import (
	"encoding/json"
	"github.com/kevinrobayna/goprod/internal"
	"github.com/stretchr/testify/assert"
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
		app := fxtest.New(t, opts(), fx.NopLogger)
		// Starts the app right away, and defers a stop when the test ends.
		defer app.RequireStart().RequireStop()

		resp, err := http.Get("http://localhost:8080/")
		if err != nil {
			panic(err.Error())
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err.Error())
		}

		var message = Message{}
		err = json.Unmarshal(body, &message)
		if err != nil {
			panic(err.Error())
		}

		assert.Equal(t, "hello", message.Message)
	})

	t.Run("Ping", func(t *testing.T) {
		t.Parallel()

		app := fxtest.New(t, opts(), fx.NopLogger)
		// Starts the app right away, and defers a stop when the test ends.
		defer app.RequireStart().RequireStop()

		resp, err := http.Get("http://localhost:8080/ping")
		if err != nil {
			panic(err.Error())
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err.Error())
		}

		var p []internal.Product
		err = json.Unmarshal(body, &p)
		if err != nil {
			panic(err.Error())
		}

		assert.Equal(t, 0, len(p))
	})
}
