package database

import (
	"context"
	"log"
	"lucienne/config"
	"time"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func ConnectDB() {
	dbURL := config.EnvVariables.DatabaseURL

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
	}

	if err := conn.Ping(ctx); err != nil {
		log.Fatalf("Erro ao dar ping no banco: %v", err)
	}

	log.Println("Conectado!")
	Conn = conn
}
