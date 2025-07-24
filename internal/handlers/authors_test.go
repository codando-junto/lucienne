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
	CreateAuthorFunc func(ctx context.Context, author *domain.Author) error
}

// Implementamos os métodos da interface AuthorRepository.
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
