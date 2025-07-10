package test_support

import "testing"

func SetupEnvVars(t *testing.T, envOverrides map[string]string) {
	t.Helper()
	t.Setenv("APP_ENV", "development")
	t.Setenv("ASSETS_PATH", "assets")
	t.Setenv("COMPILED_ASSETS_PATH", "dist/assets")

	for key, value := range envOverrides {
		t.Setenv(key, value)
	}
}
