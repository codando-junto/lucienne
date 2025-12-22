package handlers

import (
	"errors"
	"fmt"
	"log"
	"lucienne/internal/domain"
	"lucienne/internal/infra/repository"
	"lucienne/pkg/renderer"
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
	router.HandleFunc("/publishers", h.CreatePublisherHandler).Methods("POST")
	router.HandleFunc("/publishers/new", h.NewPublisherForm).Methods("GET")

	// Rota para lidar com DELETE sem ID
	router.HandleFunc("/publishers", h.MissingPublisherIDHandler).Methods("DELETE")

	// Rota para lidar com o FORM de deleção e deleção com ID
	router.HandleFunc("/publishers/{id}", h.DeletePublisherFormHandler).Methods("GET")
	router.HandleFunc("/publishers/{id}", h.DeletePublisherHandler).Methods("DELETE")
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

func (h *PublisherHandler) DeletePublisherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, _ := vars["id"]

	if idStr == "" {
		http.Error(w, "ID da editora não pode ser vazio", http.StatusBadRequest)
		return
	}

	id, err := parseID(idStr)
	if err != nil {
		http.Error(w, "ID da editora inválido", http.StatusBadRequest)
		return
	}

	err = h.repo.DeletePublisher(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrPublisherNotFound) {
			http.Error(w, "Editora não encontrada", http.StatusNotFound)
			return
		}

		log.Printf("Erro inesperado ao deletar editora: %v", err)
		http.Error(w, "Erro interno ao deletar editora", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PublisherHandler) MissingPublisherIDHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "ID da editora não fornecido", http.StatusBadRequest)
}

func (h *PublisherHandler) NewPublisherForm(w http.ResponseWriter, r *http.Request) {
	page, err := renderer.HTML.Render("publishers/new.html", nil)
	if err != nil {
		http.Error(w, "Erro ao renderizar a página", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(page)
}

func (h *PublisherHandler) DeletePublisherFormHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, _ := vars["id"]

	if idStr == "" {
		http.Error(w, "ID da editora não pode ser vazio", http.StatusBadRequest)
		return
	}

	id, err := parseID(idStr)
	if err != nil {
		http.Error(w, "ID da editora inválido", http.StatusBadRequest)
		return
	}

	data := map[string]any{
		"ID": id,
	}

	page, err := renderer.HTML.Render("publishers/delete.html", data)
	if err != nil {
		http.Error(w, "Erro ao renderizar a página", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(page)
}

// parseID converte uma string em int64, retornando um erro se a conversão falhar.
func parseID(idStr string) (int64, error) {
	var id int64
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
