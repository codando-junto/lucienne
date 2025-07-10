package config

import (
	"log"
	"os"
)

type application struct {
	RootPath    string
	Environment string
}

var Application = application{}

func (app *application) Configure(environment string) {
	app.RootPath = getRootPath()
	app.Environment = environment
	if app.Environment == "" {
		app.Environment = "development"
	}
}

func (app *application) IsDevelopment() bool {
	return app.Environment == "development"
}

func getRootPath() string {
	path, err := os.Getwd()
	for {
		_, err := os.Stat(path + "/go.mod")
		if err == nil {
			break
		}
		os.Chdir("..")
		path, _ = os.Getwd()
	}
	if err != nil {
		log.Fatalf("An error ocurred when trying to get application root path: %s", err.Error())
	}
	return path
}
