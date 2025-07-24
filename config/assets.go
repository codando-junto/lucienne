package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type assetsConfig struct {
	OriginRelativePath   string
	OriginFullPath       string
	CompiledRelativePath string
	CompiledFullPath     string
	BuildFilePath        string
	AssetsMapping        map[string]string
}

var Assets = assetsConfig{}

func (assets *assetsConfig) Configure(assetsPath string, compiledAssetsPath string, buildFilePath string) {
	assets.OriginRelativePath = removeEndingSlash(assetsPath)
	assets.CompiledRelativePath = removeEndingSlash(compiledAssetsPath)
	assets.OriginFullPath = Application.RootPath + "/" + assets.OriginRelativePath
	assets.CompiledFullPath = Application.RootPath + "/" + assets.CompiledRelativePath
	assets.BuildFilePath = Application.RootPath + "/" + buildFilePath
	assets.AssetsMapping = loadAssetsMapping()
}

func loadAssetsMapping() map[string]string {
	jsonAssets := map[string]string{}
	assetsWithPath := make(map[string]string)

	buildFile, err := os.ReadFile(Assets.BuildFilePath)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal([]byte(buildFile), &jsonAssets)
	for key, value := range jsonAssets {
		key = strings.ReplaceAll(key, Assets.OriginRelativePath+"/", "")
		value = strings.ReplaceAll(value, Assets.CompiledRelativePath+"/", "")
		assetsWithPath[key] = value
	}

	return assetsWithPath
}

func removeEndingSlash(path string) string {
	if strings.HasSuffix(path, "/") {
		return path[0 : len(path)-1]
	}
	return path
}
