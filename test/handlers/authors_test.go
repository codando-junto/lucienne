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
