package handlers

import (
	"context"
	"errors"
	"lucienne/internal/domain"
	"lucienne/internal/infra/repository"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// MockPublisherRepository é a nossa implementação falsa do repositório para testes.
type MockPublisherRepository struct {
	CreatePublisherFunc func(ctx context.Context, Publisher *domain.Publisher) error
}

// Implementamos os métodos da interface PublisherRepository.
func (m *MockPublisherRepository) CreatePublisher(ctx context.Context, Publisher *domain.Publisher) error {
	if m.CreatePublisherFunc != nil {
		return m.CreatePublisherFunc(ctx, Publisher)
	}
	return nil
}

func TestCreatePublisherHandler(t *testing.T) {
	testCases := []struct {
		name                 string
		formName             string
		mockRepo             *MockPublisherRepository
		expectedStatusCode   int
		expectedBodyContains string
	}{
		{
			name:     "deve criar uma editora com sucesso",
			formName: "Nova Editora",
			mockRepo: &MockPublisherRepository{
				CreatePublisherFunc: func(ctx context.Context, Publisher *domain.Publisher) error {
					return nil // Simula que a criação no banco foi bem-sucedida
				},
			},
			expectedStatusCode:   http.StatusCreated,
			expectedBodyContains: "Editora criada com sucesso: Nova Editora",
		},
		{
			name:     "deve retornar erro 409 ao tentar criar uma editora que já existe",
			formName: "Editora Existente",
			mockRepo: &MockPublisherRepository{
				CreatePublisherFunc: func(ctx context.Context, Publisher *domain.Publisher) error {
					return repository.ErrPublisherAlreadyExists // Simula erro de duplicidade do DB
				},
			},
			expectedStatusCode:   http.StatusConflict,
			expectedBodyContains: `Erro: A editora "Editora Existente" já está cadastrada.`,
		},
		{
			name:                 "deve retornar erro 400 se o nome estiver em branco",
			formName:             "  ",
			mockRepo:             &MockPublisherRepository{}, // O repositório não será chamado
			expectedStatusCode:   http.StatusBadRequest,
			expectedBodyContains: `O campo "name" é obrigatório`,
		},
		{
			name:     "deve retornar erro 500 se houver erro ao criar a editora",
			formName: "Editora com Falha",
			mockRepo: &MockPublisherRepository{
				CreatePublisherFunc: func(ctx context.Context, Publisher *domain.Publisher) error {
					// Simula um erro genérico do DB na criação
					return errors.New("erro de disco no banco de dados")
				},
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedBodyContains: "Erro interno ao criar editora",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Configuração do teste
			handler := NewPublisherHandler(tc.mockRepo)

			formData := url.Values{}
			formData.Set("name", tc.formName)

			req := httptest.NewRequest("POST", "/Publishers", strings.NewReader(formData.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()

			// Execução
			handler.CreatePublisherHandler(rr, req)

			// Verificação
			if status := rr.Code; status != tc.expectedStatusCode {
				t.Errorf("handler retornou status code errado: got %v want %v", status, tc.expectedStatusCode)
			}

			if !strings.Contains(rr.Body.String(), tc.expectedBodyContains) {
				t.Errorf("handler retornou corpo inesperado: got %q want to contain %q", rr.Body.String(), tc.expectedBodyContains)
			}
		})
	}
}
