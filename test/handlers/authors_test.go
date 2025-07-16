package handlers_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"lucienne/internal/handlers"
	"lucienne/internal/infra/database"

	"github.com/gorilla/mux"
)

func TestUpdateAuthor(t *testing.T) {
	database.ConnectDB()

	var id int
	err := database.Conn.QueryRow(context.Background(), "INSERT INTO authors (name) VALUES ($1) RETURNING id", "Nome Antigo").Scan(&id)
	if err != nil {
		t.Fatalf("erro ao inserir autor de teste: %v", err)
	}
	defer database.Conn.Exec(context.Background(), "DELETE FROM authors WHERE id = $1", id)

	form := url.Values{}
	form.Set("name", "Novo Nome")

	req := httptest.NewRequest("PATCH", fmt.Sprintf("/authors/%d", id), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/authors/{id}", handlers.UpdateAuthor).Methods("PATCH")
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperado status 200, recebido %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "Autor atualizado com sucesso") {
		t.Errorf("mensagem de sucesso não encontrada. Resposta: %s", rec.Body.String())
	}

	var name string
	err = database.Conn.QueryRow(context.Background(), "SELECT name FROM authors WHERE id = $1", id).Scan(&name)
	if err != nil {
		t.Fatalf("erro ao consultar autor atualizado: %v", err)
	}
	if name != "Novo Nome" {
		t.Errorf("esperado nome 'Novo Nome', recebido '%s'", name)
	}
}

func TestCreateAuthor(t *testing.T) {
	testCases := []struct {
		name         string
		formValues   map[string]string
		expectedCode int
		expectedBody string
	}{
		{
			name: "deve retornar 201 Created quando o corpo da requisição é válido",
			formValues: map[string]string{
				"name": "Teste Autor",
			},
			expectedCode: http.StatusCreated,
			expectedBody: "Autor criado com sucesso: Teste Autor",
		},
		{
			name: "deve retornar 400 Bad Request quando o nome está em branco",
			formValues: map[string]string{
				"name": "   ",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `O campo "name" é obrigatório`,
		},
		{
			name:         "deve retornar 400 Bad Request quando o formulário está vazio",
			formValues:   map[string]string{},
			expectedCode: http.StatusBadRequest,
			expectedBody: `O campo "name" é obrigatório`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			form := url.Values{}
			for key, value := range tc.formValues {
				form.Set(key, value)
			}
			requestBody := strings.NewReader(form.Encode())

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/authors", requestBody)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			handler := http.HandlerFunc(handlers.CreateAuthorHandler)
			handler.ServeHTTP(rec, req)

			if rec.Code != tc.expectedCode {
				t.Errorf("esperado status %d, recebido %d", tc.expectedCode, rec.Code)
			}
			if !strings.Contains(rec.Body.String(), tc.expectedBody) {
				t.Errorf("esperado que o corpo contivesse '%s', mas o corpo foi: '%s'", tc.expectedBody, rec.Body.String())
			}
		})
	}
}
