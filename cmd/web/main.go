package main

import (
	"github.com/kevinrobayna/goprod/internal"
	"go.uber.org/fx"
)

func opts() fx.Option {
	return fx.Options(internal.Module, WebModule)
}

func main() {
	fx.New(opts(), fx.WithLogger(internal.FxEvent)).Run()
}
