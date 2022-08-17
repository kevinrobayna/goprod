package internal

import (
	"go.uber.org/fx"
	"testing"
)

func TestApp(t *testing.T) {
	t.Parallel()

	t.Run("TestDependenciesAreSatisfied", func(t *testing.T) {
		t.Parallel()

		if err := fx.ValidateApp(Module); err != nil {
			t.Error(err)
		}
	})
}
