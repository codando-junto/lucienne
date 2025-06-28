package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Codando-Junto/ong_da_laiz/database"
	"github.com/Codando-Junto/ong_da_laiz/internal/handlers"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	database.ConnectDB()
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "9090"
	}

	r := mux.NewRouter()
	r.HandleFunc("/health", HealthHandler).Methods("GET")
	r.HandleFunc("/authors/{id}", handlers.UpdateAuthor).Methods("PATCH")

	log.Println("Rodando na porta: " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	dbURL := os.Getenv("DATABASE_URL")
	fmt.Println(dbURL)
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@postgres:5432/biblioteca?sslmode=disable"
	}

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "file://db/migrations"
	}
	m, err := migrate.New(
		migrationsPath,
		dbURL)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Iniciando migrações...")
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("Nenhuma migração pendente. Banco de dados já está atualizado.")
		} else {
			log.Fatalf("Erro ao aplicar migrações: %v", err)
		}
	} else {
		log.Println("Migrações aplicadas com sucesso.")
	}

	// Log do estado atual das migrações
	version, dirty, err := m.Version()
	if err != nil {
		log.Fatalf("Erro ao obter versão das migrações: %v", err)
	}
	log.Printf("Versão atual do banco de dados: %d, Dirty: %v", version, dirty)
}
