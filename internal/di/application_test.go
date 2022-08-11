package di

import (
	"go.uber.org/fx"
	"testing"
)

func TestApp(t *testing.T) {
	t.Parallel()

	t.Run("TestDependenciesAreSatisfied", func(t *testing.T) {
		t.Parallel()

		if err := fx.ValidateApp(ApplicationModule); err != nil {
			t.Error(err)
		}
	})
}
