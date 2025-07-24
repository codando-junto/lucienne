package repository_test

import (
	"context"
	"log"
	"lucienne/config"
	"lucienne/internal/infra/database"
	"lucienne/internal/infra/repository"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	err := godotenv.Load(".env.test")
	if err != nil {
		log.Fatalf("Erro carregando .env.test: %v", err)
	}

	config.EnvVariables.Load()
	if err != nil {
		log.Fatalf("Erro carregando config.EnvVariables: %v", err)
	}

	database.ConnectDB()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
	}

	os.Exit(m.Run())
}

// setupTestDB inicializa a conexão com o banco de dados para os testes.
func setupTestDB(t *testing.T) {
	t.Helper()
	t.Setenv("APP_ENV", "test")
	config.EnvVariables.Load()
	config.EnvVariables.DatabaseURL = config.EnvVariables.DatabaseTestURL
	database.ConnectDB()
}

func TestPostgresAuthorRepository_UpdateAuthor(t *testing.T) {
	// Só executa este teste se a flag -short não for usada
	if testing.Short() {
		t.Skip("Pulando teste de integração em modo 'short'.")
	}

	setupTestDB(t)
	ctx := context.Background()
	repo := repository.NewPostgresAuthorRepository()

	// -- Caso de Sucesso --
	t.Run("deve atualizar um autor com sucesso", func(t *testing.T) {
		// 1. Inserir um autor para o teste
		var authorID int
		originalName := "Autor Original Para Teste"
		err := database.Conn.QueryRow(ctx, "INSERT INTO authors (name) VALUES ($1) RETURNING id", originalName).Scan(&authorID)
		require.NoError(t, err, "Falha ao inserir autor para o teste de atualização")

		// Garante que o autor de teste seja removido no final
		defer func() {
			_, err := database.Conn.Exec(ctx, "DELETE FROM authors WHERE id = $1", authorID)
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
		err = database.Conn.QueryRow(ctx, "SELECT name FROM authors WHERE id = $1", authorID).Scan(&updatedName)
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
