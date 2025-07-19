package repository

import (
	"context"
	"errors"
	"lucienne/internal/domain"
	"lucienne/internal/infra/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ErrAuthorAlreadyExists é retornado quando uma tentativa de criar um autor que já existe é feita.
var ErrAuthorAlreadyExists = errors.New("author already exists")

const (
	authorExistsQuery = `SELECT 1 FROM autores WHERE nome = $1`
)

const (
	// Não precisamos retornar o ID por enquanto, então usamos um INSERT simples.
	createAuthorQuery = `INSERT INTO autores (nome) VALUES ($1)`
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

// AuthorExists verifica se um autor com o nome especificado já existe no banco de dados.
func (r *PostgresAuthorRepository) AuthorExists(ctx context.Context, authorName string) (bool, error) {
	// Usamos "SELECT 1" pois só nos importamos com a existência da linha,
	// não com os dados dela.
	var placeholder int
	err := database.Conn.QueryRow(ctx, authorExistsQuery, authorName).Scan(&placeholder)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil // Autor não existe, não é um erro.
		}
		return false, err // Erro real do banco de dados.
	}
	return true, nil // Autor existe.
}
