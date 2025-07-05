package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Codando-Junto/ong_da_laiz/internal/handlers"
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
