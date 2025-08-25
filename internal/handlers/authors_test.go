package handlers

import (
	"context"
	"errors"
	"fmt"
	"lucienne/internal/domain"
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
	CreateAuthorFunc  func(ctx context.Context, author *domain.Author) error
	UpdateAuthorFunc  func(ctx context.Context, id int, name string) error
	GetAuthorByIDFunc func(ctx context.Context, id int64) (*domain.Author, error)
	RemoveAuthorFunc  func(ctx context.Context, id int64) error
	GetAuthorsFunc    func(ctx context.Context) ([]domain.Author, error)
}

// GetAuthors implementa a interface repository.AuthorRepository.
func (m *MockAuthorRepository) GetAuthors(ctx context.Context) ([]domain.Author, error) {
	if m.GetAuthorsFunc != nil {
		return m.GetAuthorsFunc(ctx)
	}
	return nil, nil
}

// RemoveAuthor implements repository.AuthorRepository.
func (m *MockAuthorRepository) RemoveAuthor(ctx context.Context, id int64) error {
	if m.RemoveAuthorFunc != nil {
		return m.RemoveAuthorFunc(ctx, id)
	}
	return nil
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

// GetAuthorByID implementa a interface repository.AuthorRepository.
func (m *MockAuthorRepository) GetAuthorByID(ctx context.Context, id int64) (*domain.Author, error) {
	if m.GetAuthorByIDFunc != nil {
		return m.GetAuthorByIDFunc(ctx, id)
	}
	return nil, errors.New("não implementado no mock")
}

func TestNewAuthorForm(t *testing.T) {
	t.Run("deve exibir o formulário de novo autor com sucesso", func(t *testing.T) {
		// Como este handler não usa o repositório, podemos passar nil.
		handler := NewAuthorHandler(nil)
		router := mux.NewRouter()
		handler.DefineAuthors(router)

		req := httptest.NewRequest("GET", "/authors/new", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		// Verificação do Status Code
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler retornou status code errado: got %v want %v",
				status, http.StatusOK)
		}

		// Verificação do Conteúdo do Corpo
		expectedBody := "<h3>Novo Autor</h3>"
		if !strings.Contains(rr.Body.String(), expectedBody) {
			t.Errorf("handler retornou corpo inesperado: got %q want to contain %q", rr.Body.String(), expectedBody)
		}
	})
}

func TestListAuthors(t *testing.T) {
	testCases := []struct {
		name                 string
		mockRepo             *MockAuthorRepository
		expectedStatusCode   int
		expectedBodyContains []string
	}{
		{
			name: "deve listar autores com sucesso",
			mockRepo: &MockAuthorRepository{
				GetAuthorsFunc: func(ctx context.Context) ([]domain.Author, error) {
					return []domain.Author{
						{ID: 1, Name: "Autor 1"},
						{ID: 2, Name: "Autor 2"},
					}, nil
				},
			},
			expectedStatusCode:   http.StatusOK,
			expectedBodyContains: []string{"Autor 1", "Autor 2"},
		},
		{
			name: "deve retornar 500 se o repositório falhar ao listar autores",
			mockRepo: &MockAuthorRepository{
				GetAuthorsFunc: func(ctx context.Context) ([]domain.Author, error) {
					return nil, errors.New("falha de conexão com o banco")
				},
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedBodyContains: []string{"Erro interno ao listar autores"},
		},
		{
			name: "deve exibir mensagem apropriada quando não houver autores",
			mockRepo: &MockAuthorRepository{
				GetAuthorsFunc: func(ctx context.Context) ([]domain.Author, error) {
					return []domain.Author{}, nil
				},
			},
			expectedStatusCode: http.StatusOK,
			// Esta asserção depende do conteúdo de `authors/index.html`.
			expectedBodyContains: []string{"Nenhum autor encontrado"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewAuthorHandler(tc.mockRepo)
			router := mux.NewRouter()
			handler.DefineAuthors(router)

			req := httptest.NewRequest("GET", "/authors", nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatusCode {
				t.Errorf("handler retornou status code errado: got %v want %v", status, tc.expectedStatusCode)
			}

			body := rr.Body.String()
			for _, expected := range tc.expectedBodyContains {
				if !strings.Contains(body, expected) {
					t.Errorf("handler retornou corpo inesperado: got %q want to contain %q", body, expected)
				}
			}
		})
	}
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
			handler := NewAuthorHandler(tc.mockRepo)
			formData := url.Values{}
			formData.Set("name", tc.formName)

			req := httptest.NewRequest("PUT", fmt.Sprintf("/authors/%s", tc.authorID), strings.NewReader(formData.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()

			// O handler UpdateAuthor depende do mux para extrair o ID da URL.
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

func TestEditAuthorHandler(t *testing.T) {
	t.Run("deve exibir o formulário de edição com dados do autor", func(t *testing.T) {
		mockRepo := &MockAuthorRepository{
			GetAuthorByIDFunc: func(ctx context.Context, id int64) (*domain.Author, error) {
				return &domain.Author{ID: id, Name: "Autor Teste"}, nil
			},
		}
		handler := NewAuthorHandler(mockRepo)
		router := mux.NewRouter()
		handler.DefineAuthors(router)

		req := httptest.NewRequest("GET", "/authors/1/edit", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("esperava status 200, obteve %d", status)
		}
		// Verifica se o título está correto
		if !strings.Contains(rr.Body.String(), "<h2>Editar Autor</h2>") {
			t.Errorf("esperava o título '<h2>Editar Autor</h2>', mas não o encontrou no corpo: %q", rr.Body.String())
		}
		// Verifica especificamente se o campo de input está preenchido com o valor correto
		expectedInputHTML := `value="Autor Teste"`
		if !strings.Contains(rr.Body.String(), expectedInputHTML) {
			t.Errorf("esperava que o campo de input contivesse '%s', mas não encontrou no corpo: %q", expectedInputHTML, rr.Body.String())
		}
	})

	t.Run("deve retornar 404 se o autor não for encontrado", func(t *testing.T) {
		mockRepo := &MockAuthorRepository{
			GetAuthorByIDFunc: func(ctx context.Context, id int64) (*domain.Author, error) {
				return nil, repository.ErrAuthorNotFound
			},
		}
		handler := NewAuthorHandler(mockRepo)
		router := mux.NewRouter()
		handler.DefineAuthors(router)

		req := httptest.NewRequest("GET", "/authors/999/edit", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("esperava status 404, obteve %d", status)
		}
		if !strings.Contains(rr.Body.String(), "Autor não encontrado") {
			t.Errorf("esperava mensagem de autor não encontrado, obteve: %q", rr.Body.String())
		}
	})

	t.Run("deve retornar 400 se o ID for inválido", func(t *testing.T) {
		mockRepo := &MockAuthorRepository{}
		handler := NewAuthorHandler(mockRepo)
		router := mux.NewRouter()
		handler.DefineAuthors(router)

		req := httptest.NewRequest("GET", "/authors/abc/edit", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("esperava status 400, obteve %d", status)
		}
		if !strings.Contains(rr.Body.String(), "ID inválido") {
			t.Errorf("esperava mensagem de ID inválido, obteve: %q", rr.Body.String())
		}
	})

	t.Run("deve retornar 500 se o repositório falhar ao buscar autor", func(t *testing.T) {
		mockRepo := &MockAuthorRepository{
			GetAuthorByIDFunc: func(ctx context.Context, id int64) (*domain.Author, error) {
				// Simula um erro genérico do banco de dados
				return nil, errors.New("falha de conexão com o banco")
			},
		}
		handler := NewAuthorHandler(mockRepo)
		router := mux.NewRouter()
		handler.DefineAuthors(router)

		req := httptest.NewRequest("GET", "/authors/1/edit", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("esperava status 500, obteve %d", status)
		}
		if !strings.Contains(rr.Body.String(), "Erro ao buscar autor") {
			t.Errorf("esperava mensagem de erro ao buscar autor, obteve: %q", rr.Body.String())
		}
	})
}

func TestRemoveAuthorHandler(t *testing.T) {
	testCases := []struct {
		name                 string
		authorID             string
		mockRepo             *MockAuthorRepository
		expectedStatusCode   int
		expectedBodyContains string
	}{
		{
			name:     "deve remover um autor com sucesso",
			authorID: "1",
			mockRepo: &MockAuthorRepository{
				RemoveAuthorFunc: func(ctx context.Context, id int64) error {
					return nil
				},
			},
			expectedStatusCode:   http.StatusOK,
			expectedBodyContains: "Autor removido com sucesso \n",
		},
		{
			name:     "deve retornar 422 se o autor tiver livros",
			authorID: "2",
			mockRepo: &MockAuthorRepository{
				RemoveAuthorFunc: func(ctx context.Context, id int64) error {
					return repository.ErrAuthorHasBooks
				},
			},
			expectedStatusCode:   http.StatusUnprocessableEntity,
			expectedBodyContains: repository.ErrAuthorHasBooks.Error(),
		},
		{
			name:     "deve retornar 404 se o autor não for encontrado",
			authorID: "999",
			mockRepo: &MockAuthorRepository{
				RemoveAuthorFunc: func(ctx context.Context, id int64) error {
					return repository.ErrAuthorNotFound
				},
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedBodyContains: "Autor não encontrado",
		},
		{
			name:                 "deve retornar 400 se o ID for inválido",
			authorID:             "abc",
			mockRepo:             &MockAuthorRepository{},
			expectedStatusCode:   http.StatusBadRequest,
			expectedBodyContains: "ID inválido",
		},
		{
			name:     "deve retornar 500 em caso de erro genérico do repositório",
			authorID: "3",
			mockRepo: &MockAuthorRepository{
				RemoveAuthorFunc: func(ctx context.Context, id int64) error {
					return errors.New("falha de conexão com o banco")
				},
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedBodyContains: "Erro interno ao remover autor",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewAuthorHandler(tc.mockRepo)
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/authors/%s", tc.authorID), nil)
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/authors/{id}", handler.RemoveAuthor)
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
