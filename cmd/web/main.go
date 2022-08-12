package main

import (
	"github.com/kevinrobayna/goprod/internal/di"
	"go.uber.org/fx"
)

func opts() fx.Option {
	return fx.Options(di.ApplicationModule, WebModule)
}

func main() {
	fx.New(opts(), fx.WithLogger(di.FxEvent)).Run()
}
