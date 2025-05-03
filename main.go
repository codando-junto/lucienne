package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Codando-Junto/ong_da_laiz/database"
	"github.com/gorilla/mux"
)

func main() {
	database.ConnectDB()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API | teste!"))
	}).Methods("GET")
	log.Println("Rodando em http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
