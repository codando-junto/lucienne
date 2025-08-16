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

// PublisherHandler agrupa os handlers relacionados a publishers (Editoras) e suas dependências.
type PublisherHandler struct {
	repo repository.PublisherRepository
}

// NewPublisherHandler cria uma nova instância do PublisherHandler com suas dependências.
func NewPublisherHandler(repo repository.PublisherRepository) *PublisherHandler {
	return &PublisherHandler{repo: repo}
}

// DefinePublishers registra as rotas de publisher no roteador.
func (h *PublisherHandler) DefinePublishers(router *mux.Router) {
	publishersRouter := router.PathPrefix("/publishers").Subrouter()
	publishersRouter.HandleFunc("", h.CreatePublisherHandler).Methods("POST")
}

func (h *PublisherHandler) CreatePublisherHandler(w http.ResponseWriter, r *http.Request) {
	// É possível retornar qual o erro ocorreu, mas para simplificar, vamos apenas retornar um erro genérico
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

	// 2. Tenta criar o publisher no banco de dados
	publisher := &domain.Publisher{
		Name: name,
	}
	err := h.repo.CreatePublisher(r.Context(), publisher)
	if err != nil {
		// Se o repositório retornar o erro de que o publisher já existe
		//  retorna 409 Conflict.
		if errors.Is(err, repository.ErrPublisherAlreadyExists) {
			errorMessage := fmt.Sprintf("Erro: A editora %q já está cadastrada.", name)
			http.Error(w, errorMessage, http.StatusConflict)
			return
		}
		log.Printf("Erro inesperado ao criar editora: %v", err)
		http.Error(w, "Erro interno ao criar editora", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	responseMessage := fmt.Sprintf("Editora criada com sucesso: %s", name)
	w.Write([]byte(responseMessage))
}
