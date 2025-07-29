package handlers_test

import (
	"context"
	"errors"
	"fmt"
	"lucienne/internal/domain"
	"lucienne/internal/handlers"
	"lucienne/internal/infra/repository" // Importado para usar o erro customizado
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

// MockAuthorRepository é uma implementação falsa do repositório para testes unitários dos handlers.
type MockAuthorRepository struct {
	CreateAuthorFunc func(ctx context.Context, author *domain.Author) error
	UpdateAuthorFunc func(ctx context.Context, id int, name string) error
}

// CreateAuthor implementa a interface repository.AuthorRepository.
func (m *MockAuthorRepository) CreateAuthor(ctx context.Context, author *domain.Author) error {
	if m.CreateAuthorFunc != nil {
		return m.CreateAuthorFunc(ctx, author)
	}
	return nil
}

// UpdateAuthor implementa a interface repository.AuthorRepository.
func (m *MockAuthorRepository) UpdateAuthor(ctx context.Context, id int, name string) error {
	if m.UpdateAuthorFunc != nil {
		return m.UpdateAuthorFunc(ctx, id, name)
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
				CreateAuthorFunc: func(ctx context.Context, author *domain.Author) error {
					return nil // Simula que a criação no banco foi bem-sucedida
				},
			},
			expectedStatusCode:   http.StatusCreated,
			expectedBodyContains: "Autor criado com sucesso: Novo Autor",
		},
		{
			name:     "deve retornar erro 409 ao tentar criar um autor que já existe",
			formName: "Autor Existente",
			mockRepo: &MockAuthorRepository{
				CreateAuthorFunc: func(ctx context.Context, author *domain.Author) error {
					return repository.ErrAuthorAlreadyExists // Simula erro de duplicidade do DB
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
			name:     "deve retornar erro 500 se houver erro ao criar o autor",
			formName: "Autor com Falha",
			mockRepo: &MockAuthorRepository{
				CreateAuthorFunc: func(ctx context.Context, author *domain.Author) error {
					// Simula um erro genérico do DB na criação
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
			handler := handlers.NewAuthorHandler(tc.mockRepo)
			router := mux.NewRouter()
			handler.DefineAuthors(router)

			formData := url.Values{}
			formData.Set("name", tc.formName)

			req := httptest.NewRequest("POST", "/authors", strings.NewReader(formData.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()

			// Execução
			// Usamos o roteador para servir a requisição, o que é mais próximo do comportamento real.
			router.ServeHTTP(rr, req)
			// Usamos o roteador para servir a requisição, o que é mais próximo do comportamento real.
			router.ServeHTTP(rr, req)

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

func TestUpdateAuthorHandler(t *testing.T) {
	testCases := []struct {
		name                 string
		authorID             string // ID na URL, como string
		formName             string
		mockRepo             *MockAuthorRepository
		expectedStatusCode   int
		expectedBodyContains string
	}{
		{
			name:     "deve atualizar um autor com sucesso",
			authorID: "1",
			formName: "Nome Atualizado",
			mockRepo: &MockAuthorRepository{
				UpdateAuthorFunc: func(ctx context.Context, id int, name string) error {
					if id == 1 && name == "Nome Atualizado" {
						return nil // Sucesso
					}
					return errors.New("mock recebeu dados inesperados")
				},
			},
			expectedStatusCode:   http.StatusOK,
			expectedBodyContains: "Autor atualizado com sucesso",
		},
		{
			name:     "deve retornar 404 se o autor não for encontrado",
			authorID: "999",
			formName: "Nome Qualquer",
			mockRepo: &MockAuthorRepository{
				UpdateAuthorFunc: func(ctx context.Context, id int, name string) error {
					return repository.ErrAuthorNotFound // Simula erro do repositório
				},
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedBodyContains: "Autor não encontrado",
		},
		{
			name:                 "deve retornar 400 se o nome estiver em branco",
			authorID:             "1",
			formName:             "  ",
			mockRepo:             &MockAuthorRepository{}, // O repositório não será chamado
			expectedStatusCode:   http.StatusBadRequest,
			expectedBodyContains: `O campo "name" é obrigatório`,
		},
		{
			name:                 "deve retornar 400 se o ID for inválido",
			authorID:             "abc", // ID não numérico
			formName:             "Nome Válido",
			mockRepo:             &MockAuthorRepository{}, // O repositório não será chamado
			expectedStatusCode:   http.StatusBadRequest,
			expectedBodyContains: "ID inválido",
		},
		{
			name:     "deve retornar 500 em caso de erro genérico do repositório",
			authorID: "1",
			formName: "Nome Válido",
			mockRepo: &MockAuthorRepository{
				UpdateAuthorFunc: func(ctx context.Context, id int, name string) error {
					return errors.New("erro de disco no banco de dados")
				},
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedBodyContains: "Erro ao atualizar autor",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := handlers.NewAuthorHandler(tc.mockRepo)
			formData := url.Values{}
			formData.Set("name", tc.formName)
			req := httptest.NewRequest("PATCH", fmt.Sprintf("/authors/%s", tc.authorID), strings.NewReader(formData.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()

			// O handler UpdateAuthor depende do mux para extrair o ID da URL.
			// Portanto, precisamos criar um roteador para o teste.
			router := mux.NewRouter()
			router.HandleFunc("/authors/{id}", handler.UpdateAuthor)
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatusCode {
				t.Errorf("handler retornou status code errado: got %v want %v", status, tc.expectedStatusCode)
			}

			if !strings.Contains(rr.Body.String(), tc.expectedBodyContains) {
				t.Errorf("handler retornou corpo inesperado: got %q want to contain %q", rr.Body.String(), tc.expectedBodyContains)
			}
		})
	}
}
