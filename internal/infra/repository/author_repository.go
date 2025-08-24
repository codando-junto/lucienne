package repository

import (
	"context"
	"errors"
	"fmt"
	"lucienne/internal/domain"
	"lucienne/internal/infra/database"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// ErrAuthorAlreadyExists é retornado quando uma tentativa de criar um autor que já existe é feita.
	ErrAuthorAlreadyExists = errors.New("author already exists")

	// ErrAuthorNotFound é retornado quando um autor não é encontrado para uma operação.
	ErrAuthorNotFound = errors.New("autor não encontrado")

	// ErrAuthorNameCannotBeEmpty é retornado quando uma tentativa de criar ou atualizar um autor com nome vazio é feita.
	ErrAuthorNameCannotBeEmpty = errors.New("o nome do autor não pode ser vazio")

	// ErrAuthorHasBooks é retornado ao tentar remover um autor que possui livros associados.
	ErrAuthorHasBooks = errors.New("autor possui livros associados")

	// ErrSearchAuthors é retornado quando não há autores cadastrados
	ErrSearchAuthors = errors.New("error searching for authors")
)

const (
	// Não precisamos retornar o ID por enquanto, então usamos um INSERT simples.
	createAuthorQuery     = `INSERT INTO authors (name) VALUES ($1)`
	updateAuthorQuery     = `UPDATE authors SET name = $1 WHERE id = $2`
	getAuthorByIDQuery    = `SELECT id, name FROM authors WHERE id = $1`
	removeAuthorByIDQuery = `DELETE FROM authors WHERE id = $1`
	getAuthorsQuery       = `SELECT id, name FROM authors ORDER BY name ASC`
)

// AuthorRepository define a interface para as operações de autor no banco de dados.
type AuthorRepository interface {
	CreateAuthor(ctx context.Context, author *domain.Author) error
	UpdateAuthor(ctx context.Context, id int, name string) error
	GetAuthorByID(ctx context.Context, id int64) (*domain.Author, error)
	RemoveAuthor(ctx context.Context, id int64) error
	GetAuthors(ctx context.Context) ([]domain.Author, error)
}

// PostgresAuthorRepository é a implementação do AuthorRepository para o PostgreSQL.
type PostgresAuthorRepository struct {
	// No futuro, podemos adicionar o pool de conexões aqui.
}

// NewPostgresAuthorRepository cria uma nova instância do repositório.
func NewPostgresAuthorRepository() *PostgresAuthorRepository {
	return &PostgresAuthorRepository{}
}

func (r *PostgresAuthorRepository) GetAuthors(ctx context.Context) ([]domain.Author, error) {
	rows, err := database.Conn.Query(ctx, getAuthorsQuery)
	if err != nil {
		return nil, ErrSearchAuthors
	}

	// 'defer' é uma palavra especial em go. ela agenda o comando 'rows.Close()'
	// para ser executado no final da função. Isso garante que a conexão com o
	// banco seja sempre liberada
	defer rows.Close()

	authors, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Author])
	if err != nil {
		return nil, fmt.Errorf("error mapping authors: %w", err)
	}
	return authors, nil

}

// GetAuthorByID busca um autor pelo ID.
func (r *PostgresAuthorRepository) GetAuthorByID(ctx context.Context, id int64) (*domain.Author, error) {
	row := database.Conn.QueryRow(ctx, getAuthorByIDQuery, id)
	var author domain.Author
	err := row.Scan(&author.ID, &author.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAuthorNotFound
		}
		return nil, err
	}
	return &author, nil
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

// RemoveAuthor remove um autor do banco de dados, mas somente se ele não tiver livros associados.
func (r *PostgresAuthorRepository) RemoveAuthor(ctx context.Context, id int64) error {
	res, err := database.Conn.Exec(ctx, removeAuthorByIDQuery, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return ErrAuthorHasBooks
		}
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrAuthorNotFound
	}

	return nil
}
