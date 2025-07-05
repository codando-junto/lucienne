package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func DefineAuthors(router *mux.Router) {
	authorsRouter := router.PathPrefix("/authors").Subrouter()
	authorsRouter.HandleFunc("/{id}", UpdateAuthor).Methods("PATCH")
}

func UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Rota PATCH /authors/{id} OK\n"))
}
