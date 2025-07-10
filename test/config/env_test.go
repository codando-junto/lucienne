package config_test

import (
	"lucienne/config"
	"lucienne/test/test_support"
	"testing"
)

func TestAppEnvWasLoaded(t *testing.T) {
	t.Run("returns right value for AppEnv", func(t *testing.T) {
		test_support.SetupEnvVars(t, map[string]string{
			"APP_ENV": "something",
		})
		config.EnvVariables.Load()
		if config.EnvVariables.AppEnv != "something" {
			t.Errorf("Expected: something, Got: %s", config.EnvVariables.AppEnv)
		}
	})

	t.Run("panics if AppEnv is not defined", func(t *testing.T) {
		test_support.SetupEnvVars(t, map[string]string{
			"APP_ENV": "",
		})

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected: panic, Got: %s", config.EnvVariables.AppEnv)
			}
		}()

		config.EnvVariables.Load()
	})

	t.Run("returns right value for AppPort", func(t *testing.T) {
		test_support.SetupEnvVars(t, map[string]string{
			"APP_PORT": "6531",
		})
		config.EnvVariables.Load()
		if config.EnvVariables.AppPort != "6531" {
			t.Errorf("Expected: 6531, Got: %s", config.EnvVariables.AppPort)
		}
	})

	t.Run("returns default value when AppPort is not defined", func(t *testing.T) {
		test_support.SetupEnvVars(t, map[string]string{
			"APP_PORT": "",
		})
		config.EnvVariables.Load()
		if config.EnvVariables.AppPort != "9090" {
			t.Errorf("Expected: 9090, Got: %s", config.EnvVariables.AppPort)
		}
	})
}
