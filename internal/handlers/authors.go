package handlers

import (
	"errors"
	"fmt"
	"log"
	"lucienne/internal/domain"
	"lucienne/internal/infra/repository"
	"lucienne/pkg/renderer"
	"net/http"
	"strconv"
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
	router.HandleFunc("/authors/new", h.NewAuthorForm).Methods("GET")
	router.HandleFunc("/authors/{id}/edit", h.EditAuthor).Methods("GET")
	router.HandleFunc("/authors/{id}", h.UpdateAuthor).Methods("PUT", "POST")
	router.HandleFunc("/authors", h.CreateAuthorHandler).Methods("POST")
}

// EditAuthor exibe o formulário de edição de autor com dados preenchidos.
func (h *AuthorHandler) EditAuthor(w http.ResponseWriter, r *http.Request) {
	//log.Println("EditAuthor handler chamado para:", r.URL.Path)
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	author, err := h.repo.GetAuthorByID(r.Context(), int64(id))
	if errors.Is(err, repository.ErrAuthorNotFound) {
		http.Error(w, "Autor não encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Erro ao buscar autor", http.StatusInternalServerError)
		return
	}

	page, err := renderer.HTML.Render("authors/edit.html", author)
	if err != nil {
		http.Error(w, "Erro ao renderizar template", http.StatusInternalServerError)
		return
	}
	w.Write(page)
}

// NewAuthorForm exibe o formulário para criar um novo autor.
func (h *AuthorHandler) NewAuthorForm(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Rota GET /authors/new OK - o formulário de criação de autor será exibido aqui.\n"))
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

	err = h.repo.UpdateAuthor(r.Context(), id, name)
	if errors.Is(err, repository.ErrAuthorNotFound) {
		http.Error(w, "Autor não encontrado", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Erro ao atualizar autor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Autor atualizado com sucesso"))
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
