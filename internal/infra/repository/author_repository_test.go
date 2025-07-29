package repository_test

import (
	"context"
	"fmt"
	"log"
	"lucienne/config"
	"lucienne/internal/infra/database"
	"lucienne/internal/infra/repository"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tclog "github.com/testcontainers/testcontainers-go/log"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	insertQuery = "INSERT INTO authors (name) VALUES ($1) RETURNING id"
	deleteQuery = "DELETE FROM authors WHERE id = $1"
	selectQuery = "SELECT name FROM authors WHERE id = $1"
)

type testcontainerNoopLogger struct{}

func (l testcontainerNoopLogger) Printf(_ string, _ ...any) {}

// setupTestDBAndMigrate inicializa a conexão com o banco de dados para os testes, aplica as migrações e retorna uma função de limpeza.
func setupTestDBAndMigrate(t *testing.T) func() {
	t.Helper()
	t.Setenv("APP_ENV", "test")

	tclog.SetDefault(testcontainerNoopLogger{})

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		Env:          map[string]string{"POSTGRES_USER": "postgres", "POSTGRES_PASSWORD": "postgres", "POSTGRES_DB": "lucienne_test"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(10 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err, "Falha ao iniciar o container do banco de dados")

	// Obtém o host e a porta mapeada dinamicamente do container.
	endpoint, err := container.Endpoint(ctx, "")
	require.NoError(t, err, "Falha ao obter o endpoint do container")

	// Constrói a URL de conexão completa usando o endpoint dinâmico.
	databaseURL := fmt.Sprintf("postgres://postgres:postgres@%s/lucienne_test?sslmode=disable", endpoint)

	// AGORA, com o container rodando e a URL dinâmica em mãos, nós configuramos e conectamos ao banco.
	config.EnvVariables.DatabaseURL = databaseURL
	database.ConnectDB()

	// Caminho para as migrações relativo ao arquivo de teste
	migrationsPath := "file://../../../db/migrations"

	// Usa a URL dinâmica para as migrações.
	m, err := migrate.New(migrationsPath, databaseURL)
	require.NoError(t, err, "Falha ao criar instância de migração")

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err, "Falha ao aplicar migrações (Up)")
	}

	// Retorna a função de limpeza (teardown)
	return func() {
		err = m.Down()
		if err != nil && err != migrate.ErrNoChange {
			log.Printf("AVISO: Falha ao reverter migrações (Down): %v", err)
		}
		sourceErr, dbErr := m.Close()
		require.NoError(t, sourceErr)
		require.NoError(t, dbErr)
		database.Conn.Close(context.Background())

		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}
}

func TestPostgresAuthorRepository_UpdateAuthor(t *testing.T) {
	// Só executa este teste se a flag -short não for usada
	if testing.Short() {
		t.Skip("Pulando teste de integração com banco de dados em modo 'short'.")
	}

	cleanup := setupTestDBAndMigrate(t)
	defer cleanup()
	ctx := context.Background()
	repo := repository.NewPostgresAuthorRepository()

	// -- Caso de Sucesso --
	t.Run("deve atualizar um autor com sucesso", func(t *testing.T) {
		// 1. Inserir um autor para o teste
		var authorID int
		originalName := "Autor Original Para Teste"
		err := database.Conn.QueryRow(ctx, insertQuery, originalName).Scan(&authorID)
		require.NoError(t, err, "Falha ao inserir autor para o teste de atualização")

		// Garante que o autor de teste seja removido no final
		defer func() {
			_, err := database.Conn.Exec(ctx, deleteQuery, authorID)
			if err != nil {
				log.Printf("AVISO: Falha ao limpar autor de teste com ID %d: %v", authorID, err)
			}
		}()

		// 2. Chamar o método a ser testado
		newName := "Autor Atualizado"
		err = repo.UpdateAuthor(ctx, authorID, newName)
		assert.NoError(t, err)

		// 3. Verificar se o nome foi realmente atualizado no banco
		var updatedName string
		err = database.Conn.QueryRow(ctx, selectQuery, authorID).Scan(&updatedName)
		require.NoError(t, err, "Falha ao buscar autor atualizado para verificação")
		assert.Equal(t, newName, updatedName, "O nome do autor no banco de dados não corresponde ao esperado após a atualização")
	})

	// -- Caso de Falha: Autor não encontrado --
	t.Run("deve retornar ErrAuthorNotFound se o autor não existir", func(t *testing.T) {
		// Usamos um ID que é muito improvável de existir.
		nonExistentID := -999
		err := repo.UpdateAuthor(ctx, nonExistentID, "Nome Fantasma")

		// Verifica se o erro retornado é o esperado
		assert.ErrorIs(t, err, repository.ErrAuthorNotFound)
	})
}
