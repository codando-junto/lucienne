package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	port := os.Getenv("PORT")

	r := mux.NewRouter()
	r.HandleFunc("/health", HealthHandler).Methods("GET")

	fmt.Println("Servidor rodando na porta", port)
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

	requiredEnvVars := []string{"PORT"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("A variável de ambiente %s não está definida", envVar)
		}
	}
}
