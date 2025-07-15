package repository

import (
	"context"
	"errors"
	"lucienne/internal/domain"
	"lucienne/internal/infra/database"

	"github.com/jackc/pgx/v5"
)

// AuthorRepository define a interface para as operações de autor no banco de dados.
type AuthorRepository interface {
	CreateAuthor(ctx context.Context, author *domain.Author) error
	AuthorExists(ctx context.Context, authorName string) (bool, error)
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
	query := "INSERT INTO autores (nome) VALUES ($1) RETURNING id"
	err := database.Conn.QueryRow(ctx, query, author.Name).Scan(&author.ID)
	return err
}

// AuthorExists verifica se um autor com o nome especificado já existe no banco de dados.
func (r *PostgresAuthorRepository) AuthorExists(ctx context.Context, authorName string) (bool, error) {
	// Usamos "SELECT 1" pois só nos importamos com a existência da linha,
	// não com os dados dela.
	query := "SELECT 1 FROM autores WHERE nome = $1"
	var placeholder int
	err := database.Conn.QueryRow(ctx, query, authorName).Scan(&placeholder)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil // Autor não existe, não é um erro.
		}
		return false, err // Erro real do banco de dados.
	}
	return true, nil // Autor existe.
}
