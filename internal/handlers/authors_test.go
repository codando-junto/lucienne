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

// MockAuthorRepository é a nossa implementação falsa do repositório para testes.
type MockAuthorRepository struct {
	CreateAuthorFunc func(ctx context.Context, author *domain.Author) error
}

// Implementamos os métodos da interface AuthorRepository.
func (m *MockAuthorRepository) CreateAuthor(ctx context.Context, author *domain.Author) error {
	if m.CreateAuthorFunc != nil {
		return m.CreateAuthorFunc(ctx, author)
	}
	return nil
}

func TestNewAuthorForm(t *testing.T) {

	t.Run("deve retornar status 200 e o placeholder do formulário", func(t *testing.T) {
		handler := NewAuthorHandler(nil)

		// Cria uma nova requisição HTTP do tipo GET para a rota /authors/new.
		req := httptest.NewRequest("GET", "/authors/new", nil)

		rr := httptest.NewRecorder()

		// Precisamos de um roteador para despachar a requisição para o handler correto.
		// Aqui, criamos um novo roteador usando o mux.
		router := mux.NewRouter()
		// Define as rotas do AuthorHandler no roteador.
		handler.DefineAuthors(router)

		// Serve a requisição HTTP usando o roteador, que irá chamar o handler apropriado.
		router.ServeHTTP(rr, req)

		// Verifica se o código de status retornado é 200 (OK).
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler retornou status code errado: got %v want %v", status, http.StatusOK)
		}

		// Define o corpo esperado da resposta.
		expectedBody := "o formulário de criação de autor será exibido aqui"
		// Verifica se o corpo da resposta contém a string esperada.
		if !strings.Contains(rr.Body.String(), expectedBody) {
			t.Errorf("handler retornou corpo inesperado: got %q want to contain %q", rr.Body.String(), expectedBody)
		}
	})
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
			handler := NewAuthorHandler(tc.mockRepo)
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
