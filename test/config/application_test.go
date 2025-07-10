package config_test

import (
	"lucienne/config"
	"os"
	"testing"
)

func TestApplicationRootPath(t *testing.T) {
	config.Application.Configure("")
	currentPath, _ := os.Getwd()
	applicationRootPath := config.Application.RootPath
	expectedRootPath, _ := os.Getwd()
	if applicationRootPath != expectedRootPath {
		t.Errorf("Expected: %s, Got: %s", expectedRootPath, applicationRootPath)
	}
	os.Chdir(currentPath)
}

func TestApplicationEnvironment(t *testing.T) {
	t.Run("configures as development when an empty environment is sent", func(t *testing.T) {
		config.Application.Configure("")
		if config.Application.Environment != "development" {
			t.Errorf("Expected: development, Got: %s", config.Application.Environment)
		}
	})

	t.Run("configures as custom environment when environment is sent", func(t *testing.T) {
		config.Application.Configure("some_app_environment")
		if config.Application.Environment != "some_app_environment" {
			t.Errorf("Expected: some_app_environment, Got: %s", config.Application.Environment)
		}
	})
}

func TestApplicationIsDevelopment(t *testing.T) {
	t.Run("returns true when environment is development", func(t *testing.T) {
		config.Application.Configure("development")
		if !config.Application.IsDevelopment() {
			t.Errorf("Expected: true, Got: %#v", config.Application.IsDevelopment())
		}
	})

	t.Run("returns false when environment is not development", func(t *testing.T) {
		config.Application.Configure("another_env")
		if config.Application.IsDevelopment() {
			t.Errorf("Expected: false, Got: %#v", config.Application.IsDevelopment())
		}
	})
}
