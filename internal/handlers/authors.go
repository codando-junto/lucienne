package handlers

import (
	"errors"
	"fmt"
	"log"
	"lucienne/internal/domain"
	"lucienne/internal/infra/repository"
	"net/http"
	"strings"

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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Rota PATCH /authors/{id} OK\n"))
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

	// 2. Verifica se o autor já existe no banco de dados
	exists, err := h.repo.AuthorExists(r.Context(), name)
	if err != nil {
		log.Printf("Erro ao verificar a existência do autor: %v", err)
		http.Error(w, "Erro interno do servidor ao verificar autor", http.StatusInternalServerError)
		return
	}

	// 3. Se o autor já existir, retorna um erro 409 Conflict
	if exists {
		errorMessage := fmt.Sprintf("Erro: O autor '%s' já está cadastrado.", name)
		http.Error(w, errorMessage, http.StatusConflict)
		return
	}

	// 4. Se o autor não existe, prossegue com a criação
	author := &domain.Author{
		Name: name,
	}
	err = h.repo.CreateAuthor(r.Context(), author)
	if err != nil {
		// Verifica se o erro é de autor já existente, que pode ocorrer
		// apesar da verificação anterior devido a condições de corrida (race conditions).
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
