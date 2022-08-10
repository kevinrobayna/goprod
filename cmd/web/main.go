package main

import (
	"go.uber.org/fx"
	"goprod/internal/di"
)

func main() {
	fx.New(
		di.ApplicationModule, WebModule,
		fx.WithLogger(di.FxEvent),
	).Run()
}
