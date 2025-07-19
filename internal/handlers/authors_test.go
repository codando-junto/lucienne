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

// MockAuthorRepository é a nossa implementação falsa do repositório para testes.
type MockAuthorRepository struct {
	AuthorExistsFunc func(ctx context.Context, authorName string) (bool, error)
	CreateAuthorFunc func(ctx context.Context, author *domain.Author) error
}

// Implementamos os métodos da interface AuthorRepository.
func (m *MockAuthorRepository) AuthorExists(ctx context.Context, authorName string) (bool, error) {
	if m.AuthorExistsFunc != nil {
		return m.AuthorExistsFunc(ctx, authorName)
	}
	return false, nil
}

func (m *MockAuthorRepository) CreateAuthor(ctx context.Context, author *domain.Author) error {
	if m.CreateAuthorFunc != nil {
		return m.CreateAuthorFunc(ctx, author)
	}
	return nil
}

func TestCreateAuthorHandler(t *testing.T) {
	testCases := []struct {
		name                 string
		formName             string
		mockRepo             *MockAuthorRepository
		expectedStatusCode   int
		expectedBodyContains string
	}{
		{
			name:     "deve criar um autor com sucesso",
			formName: "Novo Autor",
			mockRepo: &MockAuthorRepository{
				AuthorExistsFunc: func(ctx context.Context, authorName string) (bool, error) {
					return false, nil // Simula que o autor não existe
				},
				CreateAuthorFunc: func(ctx context.Context, author *domain.Author) error {
					return nil // Simula que a criação no banco foi bem-sucedida
				},
			},
			expectedStatusCode:   http.StatusCreated,
			expectedBodyContains: "Autor criado com sucesso: Novo Autor",
		},
		{
			name:     "deve retornar erro 409 se o autor já existir",
			formName: "Autor Existente",
			mockRepo: &MockAuthorRepository{
				AuthorExistsFunc: func(ctx context.Context, authorName string) (bool, error) {
					return true, nil // Simula que o autor JÁ existe
				},
			},
			expectedStatusCode:   http.StatusConflict,
			expectedBodyContains: "Erro: O autor 'Autor Existente' já está cadastrado.",
		},
		{
			name:                 "deve retornar erro 400 se o nome estiver em branco",
			formName:             "  ",
			mockRepo:             &MockAuthorRepository{}, // O repositório não será chamado
			expectedStatusCode:   http.StatusBadRequest,
			expectedBodyContains: `O campo "name" é obrigatório`,
		},
		{
			name:     "deve retornar erro 500 se houver erro ao verificar existência",
			formName: "Qualquer Autor",
			mockRepo: &MockAuthorRepository{
				AuthorExistsFunc: func(ctx context.Context, authorName string) (bool, error) {
					return false, errors.New("erro de conexão com o banco") // Simula um erro no DB
				},
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedBodyContains: "Erro interno do servidor ao verificar autor",
		},
		{
			name:     "deve retornar erro 409 em caso de race condition na criação",
			formName: "Autor Concorrente",
			mockRepo: &MockAuthorRepository{
				AuthorExistsFunc: func(ctx context.Context, authorName string) (bool, error) {
					return false, nil // Simula que o autor não existe na primeira verificação
				},
				CreateAuthorFunc: func(ctx context.Context, author *domain.Author) error {
					// Simula que, no momento da criação, o autor já existe (race condition)
					return repository.ErrAuthorAlreadyExists
				},
			},
			expectedStatusCode:   http.StatusConflict,
			expectedBodyContains: "Erro: O autor 'Autor Concorrente' já está cadastrado.",
		},
		{
			name:     "deve retornar erro 500 se houver erro ao criar o autor",
			formName: "Autor com Falha",
			mockRepo: &MockAuthorRepository{
				AuthorExistsFunc: func(ctx context.Context, authorName string) (bool, error) {
					return false, nil // Simula que o autor não existe
				},
				CreateAuthorFunc: func(ctx context.Context, author *domain.Author) error {
					// Simula um erro genérico na criação
					return errors.New("erro de disco no banco de dados")
				},
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedBodyContains: "Erro interno ao criar autor",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Configuração do teste
			handler := NewAuthorHandler(tc.mockRepo)

			formData := url.Values{}
			formData.Set("name", tc.formName)

			req := httptest.NewRequest("POST", "/authors", strings.NewReader(formData.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()

			// Execução
			handler.CreateAuthorHandler(rr, req)

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
