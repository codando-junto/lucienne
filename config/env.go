package config

import (
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var EnvVariables = envVariables{}

type envVariables struct {
	AppEnv      string `name:"APP_ENV" required:"true"`
	AppPort     string `name:"APP_PORT" default:"9090"`
	DatabaseURL string `name:"DATABASE_URL" default:"postgres://postgres:postgres@postgres:5432/lucienne?sslmode=disable"`
}

func (envs *envVariables) Load() (envVar envVariables) {
	godotenv.Load()
	structValue := reflect.ValueOf(envs).Elem()
	structType := structValue.Type()

	var invalidVariables []string

	for i, structField := range reflect.VisibleFields(structType) {
		structTags := extractStructTags(structField)
		if setStructValue(structTags, structValue.Field(i), &invalidVariables) {
			continue
		}
	}

	if len(invalidVariables) > 0 {
		printFormattedError(invalidVariables)
	}

	return
}

type structTags struct {
	IsRequired   bool
	DefaultValue string
	EnvName      string
}

func extractStructTags(structField reflect.StructField) (structTags structTags) {
	structTags.IsRequired = false
	isRequiredStructTag := structField.Tag.Get("required")
	if isRequiredStructTag != "" {
		var err error
		structTags.IsRequired, err = strconv.ParseBool(isRequiredStructTag)
		if err != nil {
			log.Panicf("Error when trying to get a struct tag: %s", err.Error())
		}
	}

	structTags.DefaultValue = structField.Tag.Get("default")
	structTags.EnvName = structField.Tag.Get("name")

	return
}

func setStructValue(structTags structTags, structField reflect.Value, invalidVariables *[]string) bool {
	envValue := os.Getenv(structTags.EnvName)
	if structTags.IsRequired && envValue == "" {
		*invalidVariables = append(*invalidVariables, structTags.EnvName)
		return true
	}

	if envValue == "" {
		envValue = structTags.DefaultValue
	}

	structField.SetString(envValue)
	return false
}

func printFormattedError(invalidVariables []string) {
	formattedError := "- " + invalidVariables[0]
	if len(invalidVariables) > 1 {
		formattedError += "\n- " + strings.Join(invalidVariables[1:], "\n- ")
	}
	log.Panic("Error on trying to load these required variables:\n", formattedError)
}
