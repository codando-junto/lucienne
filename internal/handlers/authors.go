package handlers

import (
	"errors"
	"fmt"
	"log"
	"lucienne/internal/domain"
	"lucienne/internal/infra/repository"
	"net/http"
	"strconv"
	"strings"

	"lucienne/internal/infra/repository"

	"github.com/gorilla/mux"
)

// AuthorHandler agrupa os handlers relacionados a autores e suas dependências.
type AuthorHandler struct {
	repo repository.AuthorRepository
}

// NewAuthorHandler cria uma nova instância do AuthorHandler com suas dependências.
func NewAuthorHandler(repo repository.AuthorRepository) *AuthorHandler {
	return &AuthorHandler{repo: repo}
}

// DefineAuthors registra as rotas de autor no roteador.
func (h *AuthorHandler) DefineAuthors(router *mux.Router) {
	authorsRouter := router.PathPrefix("/authors").Subrouter()
	authorsRouter.HandleFunc("/{id}", h.UpdateAuthor).Methods("PATCH")
	authorsRouter.HandleFunc("", h.CreateAuthorHandler).Methods("POST")
}

func (h *AuthorHandler) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Erro ao ler formulário", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if strings.TrimSpace(name) == "" {
		http.Error(w, `O campo "name" é obrigatório`, http.StatusBadRequest)
		return
	}

	err = repository.UpdateAuthor(id, name)
	if errors.Is(err, repository.ErrAuthorNotFound) {
		http.Error(w, "Autor não encontrado", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Erro ao atualizar autor", http.StatusInternalServerError)
		return
	}
}

func (h *AuthorHandler) CreateAuthorHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao processar o formulário", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")

	// 1. Valida se o nome não está em branco
	if strings.TrimSpace(name) == "" {
		http.Error(w, `O campo "name" é obrigatório`, http.StatusBadRequest)
		return
	}

	// 2. Tenta criar o autor no banco de dados
	author := &domain.Author{
		Name: name,
	}
	err := h.repo.CreateAuthor(r.Context(), author)
	if err != nil {
		// Se o repositório retornar o erro de que o autor já existe
		//  retorna 409 Conflict.
		if errors.Is(err, repository.ErrAuthorAlreadyExists) {
			errorMessage := fmt.Sprintf("Erro: O autor '%s' já está cadastrado.", name)
			http.Error(w, errorMessage, http.StatusConflict)
			return
		}
		log.Printf("Erro inesperado ao criar autor: %v", err)
		http.Error(w, "Erro interno ao criar autor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	responseMessage := fmt.Sprintf("Autor criado com sucesso: %s", name)
	w.Write([]byte(responseMessage))
}
