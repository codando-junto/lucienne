package handlers

import (
	"fmt"
	"log"
	"lucienne/internal/domain"
	"lucienne/internal/infra/repository"
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

	author := &domain.Author{
		Name: name,
	}

	err := repository.CreateAuthor(r.Context(), author)
	if err != nil {
		log.Printf("Erro ao criar autor: %v", err)
		http.Error(w, "Erro ao criar autor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	responseMessage := fmt.Sprintf("Autor criado com sucesso: %s", name)
	w.Write([]byte(responseMessage))
}
