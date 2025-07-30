package repository

import (
	"context"
	"errors"
	"lucienne/internal/domain"
	"lucienne/internal/infra/database"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// ErrAuthorAlreadyExists é retornado quando uma tentativa de criar um autor que já existe é feita.
	ErrAuthorAlreadyExists = errors.New("author already exists")

	// ErrAuthorNotFound é retornado quando um autor não é encontrado para uma operação.
	ErrAuthorNotFound = errors.New("autor não encontrado")

	// ErrAuthorNameCannotBeEmpty é retornado quando uma tentativa de criar ou atualizar um autor com nome vazio é feita.
	ErrAuthorNameCannotBeEmpty = errors.New("o nome do autor não pode ser vazio")
)

const (
	// Não precisamos retornar o ID por enquanto, então usamos um INSERT simples.
	createAuthorQuery = `INSERT INTO authors (name) VALUES ($1)`
	updateAuthorQuery = `UPDATE authors SET name = $1 WHERE id = $2`
)

// AuthorRepository define a interface para as operações de autor no banco de dados.
type AuthorRepository interface {
	CreateAuthor(ctx context.Context, author *domain.Author) error
	UpdateAuthor(ctx context.Context, id int, name string) error
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

// UpdateAuthor atualiza o nome de um autor existente no banco de dados.
func (r *PostgresAuthorRepository) UpdateAuthor(ctx context.Context, id int, name string) error {
	// Adiciona validação para impedir nomes vazios.
	if strings.TrimSpace(name) == "" {
		return ErrAuthorNameCannotBeEmpty
	}

	res, err := database.Conn.Exec(ctx, updateAuthorQuery, name, id)
	if err != nil {
		// Adiciona tratamento para erro de nome duplicado
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrAuthorAlreadyExists
		}
		return err
	}

	// Verifica se alguma linha foi de fato alterada.
	rows := res.RowsAffected()
	if rows == 0 {
		return ErrAuthorNotFound
	}
	return nil
}
