package internal

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/fx"
	"path/filepath"
	"runtime"
	"time"
)

var TestModule = fx.Module("test",
	fx.Provide(
		func() BuildConfig {
			_, b, _, _ := runtime.Caller(0)

			// Root folder of this project
			Root := filepath.Join(filepath.Dir(b), "..")
			return BuildConfig{
				ConfigFile: fmt.Sprintf("%s/config.yml", Root),
				Sha:        "test",
				Date:       time.UnixDate,
			}
		},
	),
)

type DBContainer struct {
	Container testcontainers.Container
	Config    DBConfig
}

func TestWithPostgres(ctx context.Context) (DBContainer, error) {
	dsn := func(port nat.Port) string {
		cfg := config(port.Port())
		return cfg.GetDsn()
	}

	req := testcontainers.ContainerRequest{
		Image: "postgres:14.5",
		Env: map[string]string{
			"POSTGRES_USER":             "goprod",
			"POSTGRES_PASSWORD":         "",
			"POSTGRES_DB":               "goprod",
			"POSTGRES_HOST_AUTH_METHOD": "trust",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor: wait.ForSQL("5432/tcp", "postgres", dsn).
			Timeout(time.Second * 5),
	}
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return DBContainer{}, err
	}
	port, err := postgresC.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return DBContainer{}, err
	}
	return DBContainer{
		Container: postgresC,
		Config:    config(port.Port()),
	}, nil
}

func config(port string) DBConfig {
	return DBConfig{
		Host:     "localhost",
		Port:     port,
		User:     "goprod",
		Password: "",
		DbName:   "goprod",
		SslMode:  "disable",
	}
}
