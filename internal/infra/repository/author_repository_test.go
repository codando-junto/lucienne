package repository_test

import (
	"context"
	"errors"
	"lucienne/internal/infra/database"
	"lucienne/internal/infra/repository"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

const (
	insertQuery     = "INSERT INTO authors (name) VALUES ($1) RETURNING id"
	deleteQuery     = "DELETE FROM authors WHERE id = $1"
	selectQuery     = "SELECT name FROM authors WHERE id = $1"
	insertBookQuery = "INSERT INTO books (name, author_id) VALUES ($1, $2)"
)

func TestPostgresAuthorRepository_UpdateAuthor(t *testing.T) {
	setupTestDBAndMigrate(t)
	ctx := context.Background()
	repo := repository.NewPostgresAuthorRepository()

	// -- Caso de Sucesso --
	t.Run("deve atualizar um autor com sucesso", func(t *testing.T) {
		// Inserir um autor para o teste
		var authorID int
		originalName := "Autor Original Para Teste"

		err := database.Conn.QueryRow(ctx, insertQuery, originalName).Scan(&authorID)
		if err != nil {
			t.Fatalf("Falha ao inserir autor para o teste de atualização")
		}

		// Garante que o autor de teste seja removido no final
		t.Cleanup(func() {
			_, err := database.Conn.Exec(ctx, deleteQuery, authorID)
			if err != nil {
				t.Logf("AVISO: Falha ao limpar autor de teste com ID %d: %v", authorID, err)
			}
		})

		// Chamar o método a ser testado
		newName := "Autor Atualizado"
		err = repo.UpdateAuthor(ctx, authorID, newName)
		if err != nil {
			t.Errorf("esperava sucesso na atualização, mas obteve erro: %v", err)
		}

		// Verificar se o nome foi realmente atualizado no banco
		var updatedName string
		err = database.Conn.QueryRow(ctx, selectQuery, authorID).Scan(&updatedName)
		if err != nil {
			t.Fatalf("Falha ao buscar autor atualizado para verificação %s", err)
		}
		if updatedName != newName {
			t.Errorf("esperava nome '%s', mas obteve '%s'", newName, updatedName)
		}
	})

	// -- Caso de Falha: Autor não encontrado --
	t.Run("deve retornar ErrAuthorNotFound se o autor não existir", func(t *testing.T) {
		// Usamos um ID que é muito improvável de existir.
		nonExistentID := -999
		err := repo.UpdateAuthor(ctx, nonExistentID, "Nome Fantasma")

		// Verifica se o erro retornado é o esperado
		if !errors.Is(err, repository.ErrAuthorNotFound) {
			t.Errorf("esperava erro ErrAuthorNotFound, mas obteve: %v", err)
		}
	})

	// -- Caso de Falha: Nome duplicado --
	t.Run("deve retornar ErrAuthorAlreadyExists ao atualizar para um nome duplicado", func(t *testing.T) {
		// Inserir dois autores distintos
		var author1ID, author2ID int
		author1Name := "Autor Existente"
		author2Name := "Autor a ser Atualizado"

		err := database.Conn.QueryRow(ctx, insertQuery, author1Name).Scan(&author1ID)
		if err != nil {
			t.Fatalf("falha ao inserir autor: %v", err)
		}
		err = database.Conn.QueryRow(ctx, insertQuery, author2Name).Scan(&author2ID)
		if err != nil {
			t.Fatalf("falha ao inserir autor2: %v", err)
		}

		t.Cleanup(func() {
			database.Conn.Exec(ctx, deleteQuery, author1ID)
			database.Conn.Exec(ctx, deleteQuery, author2ID)
		})

		// Tentar atualizar o autor 2 com o nome do autor 1
		err = repo.UpdateAuthor(ctx, author2ID, author1Name)

		// Verificar se o erro é de autor já existente
		if !errors.Is(err, repository.ErrAuthorAlreadyExists) {
			t.Errorf("esperava erro ErrAuthorAlreadyExists, mas obteve: %v", err)
		}
	})

	// -- Caso de Falha: Nome vazio --
	t.Run("deve retornar erro ao tentar atualizar para um nome vazio", func(t *testing.T) {
		// Inserir um autor para o teste
		var authorID int
		err := database.Conn.QueryRow(ctx, insertQuery, "Autor Para Teste de Nome Vazio").Scan(&authorID)
		if err != nil {
			t.Fatalf("falha ao inserir autor: %v", err)
		}

		t.Cleanup(func() {
			database.Conn.Exec(ctx, deleteQuery, authorID)
		})

		// Tentar atualizar com um nome contendo apenas espaços
		err = repo.UpdateAuthor(ctx, authorID, "   ")
		if !errors.Is(err, repository.ErrAuthorNameCannotBeEmpty) {
			t.Errorf("esperava erro ErrAuthorNameCannotBeEmpty, mas obteve: %v", err)
		}
	})
}

func TestPostgresAuthorRepository_RemoveAuthor(t *testing.T) {
	setupTestDBAndMigrate(t)
	ctx := context.Background()
	repo := repository.NewPostgresAuthorRepository()

	// -- Caso de Sucesso: Remover autor sem livros --
	t.Run("deve remover um autor sem livros com sucesso", func(t *testing.T) {
		var authorID int64
		err := database.Conn.QueryRow(ctx, insertQuery, "Autor Para Ser Removido").Scan(&authorID)
		if err != nil {
			t.Fatalf("Falha ao inserir autor para o teste de remoção: %v", err)
		}

		err = repo.RemoveAuthor(ctx, authorID)
		if err != nil {
			t.Errorf("esperava sucesso na remoção, mas obteve erro: %v", err)
		}

		var name string
		err = database.Conn.QueryRow(ctx, selectQuery, authorID).Scan(&name)
		if !errors.Is(err, pgx.ErrNoRows) {
			t.Errorf("esperava erro pgx.ErrNoRows ao buscar autor removido, mas obteve: %v", err)
		}
	})

	// -- Caso de Falha: Autor com livros associados --
	t.Run("deve retornar ErrAuthorHasBooks se o autor tiver livros", func(t *testing.T) {
		var authorID int64
		err := database.Conn.QueryRow(ctx, insertQuery, "Autor Com Livro").Scan(&authorID)
		if err != nil {
			t.Fatalf("Falha ao inserir autor para o teste: %v", err)
		}
		t.Cleanup(func() {
			database.Conn.Exec(ctx, deleteQuery, authorID)
		})

		doneCh := make(chan struct{})
		var removeErr error
		go func() {
			removeErr = repo.RemoveAuthor(ctx, authorID)
			close(doneCh)
		}()

		go func() {
			time.Sleep(time.Second)
			_, err := database.Conn.Exec(ctx, insertBookQuery, "Livro do Autor", authorID)
			if err != nil {
				t.Errorf("Falha ao inserir livro para o teste: %v", err)
			} else {
				t.Log("inseriu livro")
			}
		}()

		<-doneCh

		if !errors.Is(removeErr, repository.ErrAuthorHasBooks) {
			t.Errorf("esperava erro ErrAuthorHasBooks, mas obteve: %v", removeErr)
		}

		var name string
		err = database.Conn.QueryRow(ctx, selectQuery, authorID).Scan(&name)
		if err != nil {
			t.Errorf("não esperava erro ao buscar autor que não deveria ter sido removido, mas obteve: %v", err)
		}
	})

	// -- Caso de Falha: Autor não encontrado --
	t.Run("deve retornar ErrAuthorNotFound se o autor não existir", func(t *testing.T) {
		nonExistentID := int64(-999)
		err := repo.RemoveAuthor(ctx, nonExistentID)

		if !errors.Is(err, repository.ErrAuthorNotFound) {
			t.Errorf("esperava erro ErrAuthorNotFound, mas obteve: %v", err)
		}
	})
}
