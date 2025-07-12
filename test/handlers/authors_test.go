package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"lucienne/internal/handlers"
)

// Testa se o endpoint PATCH /authors/{id} responde corretamente com status 200 e mensagem esperada
func TestUpdateAuthor(t *testing.T) {
	t.Run("deve responder 200 ao chamar PATCH /authors/{id}", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("PATCH", "/authors/123", nil)

		// router ou handler direto
		handler := http.HandlerFunc(handlers.UpdateAuthor)
		handler.ServeHTTP(rec, req)

		// Verifica se o status de resposta foi 200
		if rec.Code != http.StatusOK {
			t.Errorf("esperado status 200, recebido %d", rec.Code)
		}

		// Verifica se o conteúdo retornado pelo handler inclui "OK"
		if !strings.Contains(rec.Body.String(), "OK") {
			t.Errorf("esperado mensagem de sucesso, recebido: %s", rec.Body.String())
		}
	})
}

func TestCreateAuthor(t *testing.T) {
	t.Run("deve retornar 201 Created quando o corpo da requisição é válido", func(t *testing.T) {
		form := url.Values{}
		form.Set("name", "Teste Autor")
		requestBody := strings.NewReader(form.Encode())

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/authors", requestBody)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		handler := http.HandlerFunc(handlers.CreateAuthorHandler)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("esperado status %d, recebido %d", http.StatusCreated, rec.Code)
		}

		expectedBody := "Autor criado com sucesso: Teste Autor"
		if !strings.Contains(rec.Body.String(), expectedBody) {
			t.Errorf("esperado que o corpo contivesse '%s', mas o corpo foi: '%s'", expectedBody, rec.Body.String())
		}
	})

	t.Run("deve retornar 400 Bad Request quando o nome está em branco", func(t *testing.T) {
		form := url.Values{}
		form.Set("name", "   ")
		requestBody := strings.NewReader(form.Encode())

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/authors", requestBody)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		handler := http.HandlerFunc(handlers.CreateAuthorHandler)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("esperado status %d, recebido %d", http.StatusBadRequest, rec.Code)
		}

		expectedBody := `O campo "name" é obrigatório`
		if !strings.Contains(rec.Body.String(), expectedBody) {
			t.Errorf("esperado que o corpo contivesse '%s', mas o corpo foi: '%s'", expectedBody, rec.Body.String())
		}
	})
}
