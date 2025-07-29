package renderer

import (
	"bytes"
	"html/template"
	"path"
)

type templateConfig struct {
	assetsUrlPath string
	viewsDir      string
	assetsMapping map[string]string
}

// HTML gerencia a renderização de páginas html
var HTML templateConfig = templateConfig{
	assetsMapping: map[string]string{},
}

// Configure sets values to be used over the HTML rendering
// It receives the Asset URL path to be used when the page is rendered, the root views dir and the asset mapping to be parsed when some asset is provided
func (tc *templateConfig) Configure(assetsUrlPath string, viewsDir string, assetsMapping map[string]string) {
	tc.assetsUrlPath = assetsUrlPath
	tc.viewsDir = viewsDir
	tc.assetsMapping = assetsMapping
}

// Render builds an HTML with the sent data
// It receives a view to rendered, that must be on Go Template format, along with the data to be include and returns an []byte and an error.
// This function aims to encapsulate the HTML rendering implementation and provide functions to be used on the template file
func (tc templateConfig) Render(view string, data any) ([]byte, error) {
	baseFile := path.Base(view)
	tmpl, err := template.New(baseFile).Funcs(template.FuncMap{
		"assetsPath": tc.getPathToAssets,
	}).ParseFiles(path.Join(tc.viewsDir, view))

	if err != nil {
		return nil, err
	}

	buffer := &bytes.Buffer{}
	if err = tmpl.Execute(buffer, data); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (tc templateConfig) getPathToAssets(filepath string) string {
	return path.Join(HTML.assetsUrlPath, HTML.assetsMapping[filepath])
}
