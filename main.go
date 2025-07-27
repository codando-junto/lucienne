package main

import (
	"log"
	"lucienne/config"
	"lucienne/internal/handlers"
	"lucienne/internal/infra/database"
	"lucienne/internal/infra/repository"
	"lucienne/pkg/renderer"
	"net/http"
	"path"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
)

const (
	AssetsPath          = "assets"
	CompiledAssetsPath  = "public/assets"
	AssetsBuildFilePath = "public/build.json"
	ViewsPath           = "internal/views"
	AssetsServerPath    = "/assets"
	MigrationsPath      = "file://db/migrations"
	SeedsPath           = "file://db/seeds"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := renderer.HTML.Render(w, "home.html", nil); err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Ocorreu um erro ao renderizar a página"))
		}
	}).Methods("GET")
	r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	r.PathPrefix(AssetsServerPath).Handler(http.StripPrefix(AssetsServerPath, http.FileServer(http.Dir(CompiledAssetsPath))))

	// Injeção de Dependência
	authorRepo := repository.NewPostgresAuthorRepository()
	authorHandler := handlers.NewAuthorHandler(authorRepo)

	handlers.ReturnHealth(r)
	authorHandler.DefineAuthors(r)

	log.Println("Rodando na porta: " + config.EnvVariables.AppPort)
	log.Fatal(http.ListenAndServe(":"+config.EnvVariables.AppPort, r))
}

func init() {
	config.EnvVariables.Load()
	config.Application.Configure(config.EnvVariables.AppEnv)
	config.Assets.Configure(AssetsPath, CompiledAssetsPath, AssetsBuildFilePath)
	renderer.HTML.Configure(AssetsServerPath, path.Join(config.Application.RootPath, ViewsPath), config.Assets.AssetsMapping)
	database.ConnectDB()

	m, err := migrate.New(
		MigrationsPath,
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

	// Fechar a instância de migração para liberar a conexão com o banco de dados.
	sourceErr, dbErr := m.Close()
	if sourceErr != nil {
		log.Fatalf("Erro ao fechar o source da migração: %v", sourceErr)
	}
	if dbErr != nil {
		log.Fatalf("Erro ao fechar a conexão do banco de dados da migração: %v", dbErr)
	}

	if config.Application.IsDevelopment() {
		log.Println("Ambiente de desenvolvimento detectado. Aplicando seed...")
		seed, err := migrate.New(
			SeedsPath,
			config.EnvVariables.DatabaseURL)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Iniciando seed...")
		if err := seed.Up(); err != nil {
			if err == migrate.ErrNoChange {
				log.Println("Nenhum seed pendente. Banco de dados já está atualizado.")
			} else {
				log.Fatalf("Erro ao aplicar seeds: %v", err)
			}
		} else {
			log.Println("Seed aplicadas com sucesso.")
		}

		sourceErr, dbErr := seed.Close()
		if sourceErr != nil {
			log.Fatalf("Erro ao fechar o source do seed: %v", sourceErr)
		}
		if dbErr != nil {
			log.Fatalf("Erro ao fechar a conexão do banco de dados do seed: %v", dbErr)
		}
	}

}
