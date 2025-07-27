package renderer

import (
	"html/template"
	"io"
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

func (tc templateConfig) Render(writer io.Writer, view string, data any) error {
	baseFile := path.Base(view)
	tmpl, err := template.New(baseFile).Funcs(template.FuncMap{
		"assetsPath": tc.getPathToAssets,
	}).ParseFiles(path.Join(tc.viewsDir, view))

	if err != nil {
		return err
	}

	return tmpl.Execute(writer, data)
}

func (tc templateConfig) getPathToAssets(filepath string) string {
	return path.Join(HTML.assetsUrlPath, HTML.assetsMapping[filepath])
}
