package handlers

import (
	"net/http"
)

func CreateAuthorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Autor criado com sucesso"))
}
