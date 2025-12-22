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

	"github.com/gorilla/mux"
)

// MockPublisherRepository é a nossa implementação falsa do repositório para testes.
type MockPublisherRepository struct {
	CreatePublisherFunc func(ctx context.Context, Publisher *domain.Publisher) error
	DeletePublisherFunc func(ctx context.Context, id int64) error
}

// Implementamos os métodos da interface PublisherRepository.
func (m *MockPublisherRepository) CreatePublisher(ctx context.Context, Publisher *domain.Publisher) error {
	if m.CreatePublisherFunc != nil {
		return m.CreatePublisherFunc(ctx, Publisher)
	}
	return nil
}

func (m *MockPublisherRepository) DeletePublisher(ctx context.Context, id int64) error {
	if m.DeletePublisherFunc != nil {
		return m.DeletePublisherFunc(ctx, id)
	}
	return nil
}

func TestNewPublisherForm(t *testing.T) {
	handler := NewPublisherHandler(nil)
	router := mux.NewRouter()
	handler.DefinePublishers(router)

	req := httptest.NewRequest("GET", "/publishers/new", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler retornou status code errado: got %v want %v", status, http.StatusOK)
	}

	expectedBody := `<h1>Criar Editora</h1>
    <form action="/publishers" method="post">
        <label for="name">Nome</label>
        <input type="text" id="name" name="name" required>
        <button type="submit">Criar</button>
    </form>`

	if !strings.Contains(rr.Body.String(), expectedBody) {
		t.Errorf("handler retornou corpo inesperado: got %q want %q", rr.Body.String(), expectedBody)
	}
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

func TestDeletePublisherForm(t *testing.T) {
	handler := NewPublisherHandler(nil)
	router := mux.NewRouter()
	handler.DefinePublishers(router)

	req := httptest.NewRequest("GET", "/publishers/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler retornou status code errado: got %v want %v", status, http.StatusOK)
	}

	expectedBody := `<h1>Deletar Editora</h1>
    <form action="/publishers/1" method="delete">
        <input type="hidden" name="_method" value="delete">
        <p>Tem certeza que deseja deletar a editora?</p>
        <button type="submit">Deletar</button>
    </form>`

	if !strings.Contains(rr.Body.String(), expectedBody) {
		t.Errorf("handler retornou corpo inesperado: got %q want %q", rr.Body.String(), expectedBody)
	}
}

func TestMissingPublisherIDHandler(t *testing.T) {
	handler := NewPublisherHandler(nil)
	router := mux.NewRouter()
	handler.DefinePublishers(router)

	req := httptest.NewRequest("DELETE", "/publishers", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler retornou status code errado: got %v want %v", status, http.StatusBadRequest)
	}

	expectedBody := "ID da editora não fornecido"

	if !strings.Contains(rr.Body.String(), expectedBody) {
		t.Errorf("handler retornou corpo inesperado: got %q want %q", rr.Body.String(), expectedBody)
	}
}

func TestDeletePublisherHandler(t *testing.T) {
	testCases := []struct {
		name                 string
		url                  string
		mockRepo             *MockPublisherRepository
		expectedStatusCode   int
		expectedBodyContains string
	}{
		{
			name: "deve deletar uma editora com sucesso",
			url:  "/publishers/1",
			mockRepo: &MockPublisherRepository{
				DeletePublisherFunc: func(ctx context.Context, id int64) error {
					return nil // Simula que a deleção no banco foi bem-sucedida
				},
			},
			expectedStatusCode:   http.StatusNoContent,
			expectedBodyContains: "",
		},
		{
			name: "deve retornar erro 404 se a editora não for encontrada",
			url:  "/publishers/999",
			mockRepo: &MockPublisherRepository{
				DeletePublisherFunc: func(ctx context.Context, id int64) error {
					return repository.ErrPublisherNotFound // Simula que a editora não foi encontrada
				},
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedBodyContains: "Editora não encontrada",
		},
		{
			name: "deve retornar erro 500 se houver erro ao deletar a editora",
			url:  "/publishers/2",
			mockRepo: &MockPublisherRepository{
				DeletePublisherFunc: func(ctx context.Context, id int64) error {
					// Simula um erro genérico do DB na deleção
					return errors.New("erro de disco no banco de dados")
				},
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedBodyContains: "Erro interno ao deletar editora",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Configuração do teste
			handler := NewPublisherHandler(tc.mockRepo)

			req := httptest.NewRequest("DELETE", tc.url, nil)
			rr := httptest.NewRecorder()

			// Execução
			router := mux.NewRouter()
			handler.DefinePublishers(router)
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
