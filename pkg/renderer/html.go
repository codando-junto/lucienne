package renderer

import (
	"fmt"
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
	fmt.Println(viewsDir)
	tc.assetsUrlPath = assetsUrlPath
	tc.viewsDir = viewsDir
	tc.assetsMapping = assetsMapping
}

func (tc templateConfig) Render(writer io.Writer, view string, data any) {
	baseFile := path.Base(view)
	tmpl, err := template.New(baseFile).Funcs(template.FuncMap{
		"assetsPath": tc.getPathToAssets,
	}).ParseFiles(path.Join(tc.viewsDir, view))

	if err != nil {
		fmt.Println("Error on rendering HTML:\n" + err.Error())
	}

	tmpl.Execute(writer, data)
}

func (tc templateConfig) getPathToAssets(filepath string) string {
	return path.Join(HTML.assetsUrlPath, HTML.assetsMapping[filepath])
}
