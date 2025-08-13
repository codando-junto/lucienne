package handlers

import (
	"lucienne/config"
	"lucienne/pkg/renderer"
	"os"
	"path"
	"testing"
)

func TestMain(m *testing.M) {
	// Configura as dependências necessárias para os handlers deste pacote antes de rodar os testes.
	// Isso garante que os testes sejam autocontidos e não dependam da inicialização do pacote main.
	config.Application.Configure("test")
	viewsPath := path.Join(config.Application.RootPath, "internal/views")
	renderer.HTML.Configure("", viewsPath, nil)

	// Roda todos os testes do pacote
	exitCode := m.Run()

	os.Exit(exitCode)
}
