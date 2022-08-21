package main

import (
	"github.com/kevinrobayna/goprod/internal"
	"go.uber.org/fx"
	"testing"
)

func TestWebApp(t *testing.T) {
	t.Parallel()

	t.Run("TestDependenciesAreSatisfied", func(t *testing.T) {
		t.Parallel()

		if err := fx.ValidateApp(internal.TestModule, opts(), fx.NopLogger); err != nil {
			t.Error(err)
		}
	})
}
