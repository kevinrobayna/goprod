package internal

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type DBContainer struct {
	Container testcontainers.Container
	Config    DBConfig
}

func TestWithPostgres(ctx context.Context) (DBContainer, error) {
	req := testcontainers.ContainerRequest{
		Image: "postgres:14.5",
		Env: map[string]string{
			"POSTGRES_USER":             "goprod",
			"POSTGRES_PASSWORD":         "",
			"POSTGRES_DB":               "goprod",
			"POSTGRES_HOST_AUTH_METHOD": "trust",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
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
