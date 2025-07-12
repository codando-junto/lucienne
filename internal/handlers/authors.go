package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type CreateAuthorRequest struct {
	Name string `json:"name"`
}

func DefineAuthors(router *mux.Router) {
	authorsRouter := router.PathPrefix("/authors").Subrouter()
	authorsRouter.HandleFunc("/{id}", UpdateAuthor).Methods("PATCH")
	authorsRouter.HandleFunc("", CreateAuthorHandler).Methods("POST")
}

func UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Rota PATCH /authors/{id} OK\n"))
}

func CreateAuthorHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao processar o formulário", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")

	if strings.TrimSpace(name) == "" {

		http.Error(w, `O campo "name" é obrigatório`, http.StatusBadRequest)
		return
	}

	// Aqui viria a lógica para salvar o autor no banco de dados...
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("Autor criado com sucesso: " + name))

}
