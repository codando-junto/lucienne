package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func ReturnHealth(router *mux.Router) {
	router.HandleFunc("/health", HealthHandler).Methods("GET")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("HEALTH OK"))
}
