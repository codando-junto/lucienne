package renderer

import (
	"os"
	"path"
	"strings"
	"testing"
)

func TestRederingHTML(t *testing.T) {
	tempDir := t.TempDir()
	HTML.Configure("/assets", tempDir, map[string]string{"some_path/something.asset": "some_path/other_path/random.asset"})
	setupHTMLFile(t, path.Join(tempDir, "test.html"))

	t.Run("render inner variables", func(t *testing.T) {
		page, _ := HTML.Render("test.html", map[string]string{"TestContent": "some content"})
		if !strings.Contains(string(page), "some content") {
			t.Error("Expected: contains rendered value \"some content\", got: nothing")
		}
	})

	t.Run("render asset path", func(t *testing.T) {
		page, _ := HTML.Render("test.html", map[string]string{"TestContent": "some content"})
		if !strings.Contains(string(page), "<script src=/assets/some_path/other_path/random.asset></script>") {
			t.Error("Expected: contains rendered asset path \"/assets/some_path/other_path/random.asset\", got: nothing")
		}
	})

	t.Run("returns an error when file does not exist", func(t *testing.T) {
		_, err := HTML.Render("missing.html", map[string]string{"TestContent": "some content"})
		if err == nil {
			t.Error("Expected: some error, got: nothing")
		}
	})
}

func setupHTMLFile(t testing.TB, filePath string) {
	t.Helper()

	htmlContent := []byte(`
		<html>
			<head>
				<title>Testing</title>
			</head>
			<body>
				{{ .TestContent }}
				<script src={{ assetsPath "some_path/something.asset" }}></script>
			</body>
		</html>
	`)

	if err := os.WriteFile(filePath, htmlContent, 0644); err != nil {
		t.Fatalf("failed to write test HTML file: %v'", err)
	}

	t.Cleanup(func() {
		os.Remove(filePath)
	})
}
