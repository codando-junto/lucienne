package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"lucienne/internal/infra/repository"

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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Autor atualizado com sucesso\n"))
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
