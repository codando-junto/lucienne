package repository

import (
	"context"
	"lucienne/internal/domain"
	"lucienne/internal/infra/database"
)

// CreateAuthor insere um novo autor no banco de dados.
func CreateAuthor(ctx context.Context, author *domain.Author) error {
	query := "INSERT INTO autores (nome) VALUES ($1) RETURNING id"

	// Executa a query e escaneia o ID retornado para dentro da struct do autor.
	err := database.Conn.QueryRow(ctx, query, author.Name).Scan(&author.ID)

	return err
}
