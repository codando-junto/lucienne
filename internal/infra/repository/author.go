package repository

import (
	"context"
	"errors"
	"lucienne/internal/infra/database"
)

var ErrAuthorNotFound = errors.New("autor não encontrado")

func UpdateAuthor(id int, name string) error {
	res, err := database.Conn.Exec(context.Background(), "UPDATE authors SET name = $1 WHERE id = $2", name, id)
	if err != nil {
		return err
	}
	rows := res.RowsAffected()
	if rows == 0 {
		return ErrAuthorNotFound
	}
	return nil
}
