package repository

import (
	"context"
	"errors"
	"lucienne/internal/domain"
	"lucienne/internal/infra/database"

	"github.com/jackc/pgx/v5/pgconn"
)

// ErrAuthorAlreadyExists é retornado quando uma tentativa de criar um autor que já existe é feita.
var ErrAuthorAlreadyExists = errors.New("author already exists")

const (
	// Não precisamos retornar o ID por enquanto, então usamos um INSERT simples.
	createAuthorQuery = `INSERT INTO authors (name) VALUES ($1)`
)

// AuthorRepository define a interface para as operações de autor no banco de dados.
type AuthorRepository interface {
	CreateAuthor(ctx context.Context, author *domain.Author) error
}

// PostgresAuthorRepository é a implementação do AuthorRepository para o PostgreSQL.
type PostgresAuthorRepository struct {
	// No futuro, podemos adicionar o pool de conexões aqui.
}

// NewPostgresAuthorRepository cria uma nova instância do repositório.
func NewPostgresAuthorRepository() *PostgresAuthorRepository {
	return &PostgresAuthorRepository{}
}

// CreateAuthor insere um novo autor no banco de dados.
func (r *PostgresAuthorRepository) CreateAuthor(ctx context.Context, author *domain.Author) error {
	_, err := database.Conn.Exec(ctx, createAuthorQuery, author.Name)
	if err != nil {
		// Verifica se o erro é uma violação de chave única (unique_violation).
		// O código '23505' é o código de erro padrão do PostgreSQL para isso.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrAuthorAlreadyExists
		}
		// Se for outro tipo de erro, retorna o erro original.
		return err
	}
	return nil
}
