package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var Conn *pgx.Conn

func ConnectDB() {
	if err := godotenv.Load(); err != nil {
		log.Println("Erro ao carregar o arquivo .env")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
	}

	if err := conn.Ping(ctx); err != nil {
		log.Fatalf("Erro ao dar ping no banco: %v", err)
	}

	log.Println("Conectado!")
	Conn = conn
}
