package handlers

import (
	"net/http"
)

func UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`Rota PATCH /authors/{id} OK`))
}
