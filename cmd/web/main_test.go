package main

import (
	"go.uber.org/fx"
	"testing"
)

func TestWebApp(t *testing.T) {
	t.Parallel()

	t.Run("TestDependenciesAreSatisfied", func(t *testing.T) {
		t.Parallel()

		if err := fx.ValidateApp(opts(), fx.NopLogger); err != nil {
			t.Error(err)
		}
	})
}
