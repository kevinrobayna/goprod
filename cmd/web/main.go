package main

import (
	"go.uber.org/fx"
	"goprod/internal/di"
)

func opts() fx.Option {
	return fx.Options(di.ApplicationModule, WebModule)
}

func main() {
	fx.New(opts(), fx.WithLogger(di.FxEvent)).Run()
}
