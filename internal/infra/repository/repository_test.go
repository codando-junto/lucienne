package repository_test

import (
	"context"
	"fmt"
	"lucienne/config"
	"lucienne/internal/infra/database"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	tclog "github.com/testcontainers/testcontainers-go/log"
	"github.com/testcontainers/testcontainers-go/wait"
)

type testcontainerNoopLogger struct{}

func (l testcontainerNoopLogger) Printf(_ string, _ ...any) {}

// setupTestDBAndMigrate inicializa a conexão com o banco de dados para os testes, aplica as migrações e retorna uma função de limpeza.
func setupTestDBAndMigrate(t testing.TB) {
	t.Helper()
	t.Setenv("APP_ENV", "test")

	tclog.SetDefault(testcontainerNoopLogger{})

	ctx := t.Context()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		Env:          map[string]string{"POSTGRES_USER": "postgres", "POSTGRES_PASSWORD": "postgres", "POSTGRES_DB": "lucienne_test"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor: wait.ForSQL("5432/tcp", "pgx/v5", func(host string, port nat.Port) string {
			return fmt.Sprintf("postgres://postgres:postgres@%s:%s/lucienne_test?sslmode=disable", host, port.Port())
		}).WithStartupTimeout(10 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Falha ao iniciar o container do banco de dados %s", err)
	}

	// Obtém o host e a porta mapeada dinamicamente do container.
	endpoint, err := container.Endpoint(ctx, "")
	if err != nil {
		t.Fatalf("Falha ao obter o endpoint do container %s", err)
	}
	// Constrói a URL de conexão completa usando o endpoint dinâmico.
	databaseURL := fmt.Sprintf("postgres://postgres:postgres@%s/lucienne_test?sslmode=disable", endpoint)

	// AGORA, com o container rodando e a URL dinâmica em mãos, nós configuramos e conectamos ao banco.
	config.EnvVariables.DatabaseURL = databaseURL
	database.ConnectDB()

	// Caminho para as migrações relativo ao arquivo de teste
	migrationsPath := "file://../../../db/migrations"

	// Usa a URL dinâmica para as migrações.
	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		t.Fatalf("Falha ao criar instância de migração: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Falha ao aplicar migrações (Up): %v", err)
	}

	// Retorna a função de limpeza (teardown)
	t.Cleanup(func() {
		err = m.Down()
		if err != nil && err != migrate.ErrNoChange {
			t.Logf("AVISO: Falha ao reverter migrações (Down): %v", err)
		}

		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			t.Logf("AVISO: Falha ao fechar source da migração: %v", sourceErr)
		}
		if dbErr != nil {
			t.Logf("AVISO: Falha ao fechar DB da migração: %v", dbErr)
		}

		database.Conn.Close(ctx)

		if err := container.Terminate(context.Background()); err != nil {
			t.Logf("AVISO: Falha ao terminar o container: %s", err)
		}
	})
}
