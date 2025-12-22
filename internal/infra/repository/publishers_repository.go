package repository

import (
	"context"
	"errors"
	"lucienne/internal/domain"
	"lucienne/internal/infra/database"

	"github.com/jackc/pgx/v5/pgconn"
)

// ErrPublisherAlreadyExists é retornado quando uma tentativa de criar um autor que já existe é feita.
var ErrPublisherAlreadyExists = errors.New("Publisher already exists")

// ErrPublisherNotFound é retornado quando um autor não é encontrado no banco de dados.
var ErrPublisherNotFound = errors.New("Publisher not found")

const (
	// Não precisamos retornar o ID por enquanto, então usamos um INSERT simples.
	createPublisherQuery = `INSERT INTO publishers (name) VALUES ($1)`

	deletePublisherQuery = `DELETE FROM publishers WHERE id = $1`
)

// PublisherRepository define a interface para as operações de publisher no banco de dados.
type PublisherRepository interface {
	CreatePublisher(ctx context.Context, Publisher *domain.Publisher) error
	DeletePublisher(ctx context.Context, id int64) error
}

// PostgresPublisherRepository é a implementação do PublisherRepository para o PostgreSQL.
type PostgresPublisherRepository struct {
	// No futuro, podemos adicionar o pool de conexões aqui.
}

// NewPostgresPublisherRepository cria uma nova instância do repositório.
func NewPostgresPublisherRepository() *PostgresPublisherRepository {
	return &PostgresPublisherRepository{}
}

// CreatePublisher insere um novo publisher no banco de dados.
func (r *PostgresPublisherRepository) CreatePublisher(ctx context.Context, Publisher *domain.Publisher) error {
	_, err := database.Conn.Exec(ctx, createPublisherQuery, Publisher.Name)
	if err != nil {
		// Verifica se o erro é uma violação de chave única (unique_violation).
		// O código '23505' é o código de erro padrão do PostgreSQL para isso.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrPublisherAlreadyExists
		}
		// Se for outro tipo de erro, retorna o erro original.
		return err
	}
	return nil
}

// DeletePublisher remove um publisher do banco de dados pelo ID.
func (r *PostgresPublisherRepository) DeletePublisher(ctx context.Context, id int64) error {
	cmdTag, err := database.Conn.Exec(ctx, deletePublisherQuery, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrPublisherNotFound
	}

	return nil
}

// TODO:
// Para facilitar, podemos criar uma função que faz a busca por ID,
// que pode ser utilizada nas outras  operaçoes de CRUD.
// Ex.
// fetchPublisherByID(ctx context.Context, id int64) (*domain.Publisher, error)
