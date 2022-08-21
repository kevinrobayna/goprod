package main

import (
	"github.com/kevinrobayna/goprod/internal"
	"github.com/kevinrobayna/goprod/internal/web"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"log"
	"os"
	"runtime"
	"sort"
)

var (
	AppName   = "goprod"
	Version   = "unknown"
	Date      = "unknown"
	GoVersion = runtime.Version()
)

func opts() fx.Option {
	return fx.Options(web.Module)
}

func main() {
	app := &cli.App{
		EnableBashCompletion: true,
		Name:                 AppName,
		Usage:                "A simple web service",
		Version:              Version,
		Authors: []*cli.Author{
			{
				Name:  "Kevin Robayna",
				Email: "me@kevinrobayna.com",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Value:       "config.yml",
				Usage:       "Config file from which to load configuration",
				DefaultText: "config.yml",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "web",
				Usage: "start the web server",
				Action: func(cCtx *cli.Context) error {
					fx.New(
						opts(),
						fx.Provide(
							func() internal.BuildConfig {
								return internal.BuildConfig{
									ConfigFile: cCtx.String("config"),
									Sha:        Version,
									Date:       Date,
								}
							},
						),
						fx.WithLogger(internal.FxEvent),
					).Run()
					return nil
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
