package main

import (
	"log"
	"lucienne/config"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
)

const (
	MIGRATIONS_PATH = "file://db/migrations"
	SEEDS_PATH      = "file://db/seeds"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/health", HealthHandler).Methods("GET")

	log.Println("Rodando na porta: " + config.EnvVariables.AppPort)
	log.Fatal(http.ListenAndServe(":"+config.EnvVariables.AppPort, r))
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func init() {
	config.EnvVariables.Load()

	m, err := migrate.New(
		MIGRATIONS_PATH,
		config.EnvVariables.DatabaseURL,
	)
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

	if config.EnvVariables.AppEnv == "development" {
		log.Println("Ambiente de desenvolvimento detectado. Aplicando seed...")
		seed, err := migrate.New(
			SEEDS_PATH,
			config.EnvVariables.DatabaseURL)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Iniciando seed...")
		if err := seed.Up(); err != nil {
			if err == migrate.ErrNoChange {
				log.Println("Nenhum seed pendente. Banco de dados já está atualizado.")
			} else {
				log.Fatalf("Erro ao aplicar migrações: %v", err)
			}
		} else {
			log.Println("Seed aplicadas com sucesso.")
		}
	}

}
