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

var HTML templateConfig = templateConfig{
	assetsMapping: map[string]string{},
}

func (tc *templateConfig) Configure(assetsUrlPath string, viewsDir string, assetsMapping map[string]string) {
	tc.assetsUrlPath = assetsUrlPath
	tc.viewsDir = viewsDir
	tc.assetsMapping = assetsMapping
}

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
