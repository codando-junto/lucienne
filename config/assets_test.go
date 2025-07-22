package config

import (
	"maps"
	"os"
	"testing"
)

func TestAssetsConfigure(t *testing.T) {
	Application.Configure("test")
	rootDir, _ := os.Getwd()

	t.Run("set value for OriginRelativePath", func(t *testing.T) {
		Assets.Configure("my_origin_assets", "test/config/test_compiled_path/", "test/config/test_compiled_path/test-build.json")
		expectedResult := "my_origin_assets"
		if expectedResult != Assets.OriginRelativePath {
			t.Errorf("Expected: %s, Got: %s", expectedResult, Assets.OriginRelativePath)
		}
	})

	t.Run("removes ending slash from OriginRelativePath", func(t *testing.T) {
		Assets.Configure("my_origin_assets/", "test/config/test_compiled_path/", "test/config/test_compiled_path/test-build.json")
		expectedResult := "my_origin_assets"
		if expectedResult != Assets.OriginRelativePath {
			t.Errorf("Expected: %s, Got: %s", expectedResult, Assets.OriginRelativePath)
		}
	})

	t.Run("sets right value for CompiledRelativePath", func(t *testing.T) {
		Assets.Configure("my_origin_assets/", "test/config/test_compiled_path", "test/config/test_compiled_path/test-build.json")
		expectedResult := "test/config/test_compiled_path"
		if expectedResult != Assets.CompiledRelativePath {
			t.Errorf("Expected: %s, Got: %s", expectedResult, Assets.CompiledRelativePath)
		}
	})

	t.Run("removes ending slash CompiledRelativePath", func(t *testing.T) {
		Assets.Configure("my_origin_assets/", "test/config/test_compiled_path", "test/config/test_compiled_path/test-build.json")
		expectedResult := "test/config/test_compiled_path"
		if expectedResult != Assets.CompiledRelativePath {
			t.Errorf("Expected: %s, Got: %s", expectedResult, Assets.CompiledRelativePath)
		}
	})

	t.Run("sets right value for BuildFile", func(t *testing.T) {
		Assets.Configure("my_origin_assets/", "test/config/test_compiled_path/", "test/config/test_compiled_path/test-build.json")
		expectedResult := rootDir + "/test/config/test_compiled_path/test-build.json"
		if expectedResult != Assets.BuildFilePath {
			t.Errorf("Expected: %s, Got: %s", expectedResult, Assets.BuildFilePath)
		}
	})

	t.Run("sets right OriginFullPath", func(t *testing.T) {
		Assets.Configure("my_origin_assets/", "test/config/test_compiled_path/", "test/config/test_compiled_path/test-build.json")
		expectedResult := rootDir + "/my_origin_assets"
		if expectedResult != Assets.OriginFullPath {
			t.Errorf("Expected: %s, Got: %s", expectedResult, Assets.OriginFullPath)
		}
	})

	t.Run("removes ending slash from OriginFullPath", func(t *testing.T) {
		Assets.Configure("my_origin_assets/", "test/config/test_compiled_path/", "test/config/test_compiled_path/test-build.json")
		expectedResult := rootDir + "/my_origin_assets"
		if expectedResult != Assets.OriginFullPath {
			t.Errorf("Expected: %s, Got: %s", expectedResult, Assets.OriginFullPath)
		}
	})

	t.Run("sets right CompiledFullPath", func(t *testing.T) {
		Assets.Configure("my_origin_assets/", "test/config/test_compiled_path/", "test/config/test_compiled_path/test-build.json")
		expectedResult := rootDir + "/test/config/test_compiled_path"
		if expectedResult != Assets.CompiledFullPath {
			t.Errorf("Expected: %s, Got: %s", expectedResult, Assets.CompiledFullPath)
		}
	})

	t.Run("removes ending slash from CompiledFullPath", func(t *testing.T) {
		Assets.Configure("my_origin_assets/", "test/config/test_compiled_path/", "test/config/test_compiled_path/test-build.json")
		expectedResult := rootDir + "/test/config/test_compiled_path"
		if expectedResult != Assets.CompiledFullPath {
			t.Errorf("Expected: %s, Got: %s", expectedResult, Assets.CompiledFullPath)
		}
	})

	t.Run("loads build json mapping file", func(t *testing.T) {
		Assets.Configure("my_origin_assets/", "test/config/test_compiled_path/", "test/config/test_compiled_path/test-build.json")
		expectedResult := map[string]string{
			"some_dir_1/some-file.js":      "some_dir_1/first-file-AJFJEO.js",
			"some_dir_2/another-file.scss": "some_dir_2/another-css-AJITD2.css",
			"random-file.jpg":              "random-file-GDSJOQR.jpg",
		}
		if !maps.Equal(expectedResult, Assets.AssetsMapping) {
			t.Errorf("Expected: %#v, Got: %#v", expectedResult, Assets.AssetsMapping)
		}
	})
}
