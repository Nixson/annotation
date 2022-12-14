package annotation

import (
	"embed"
	"encoding/json"
	"github.com/Nixson/environment"
)

type Element struct {
	Url        string            `json:"url"`
	Type       string            `json:"type"`
	StructName string            `json:"structName"`
	Parameters map[string]string `json:"parameters"`
	Children   []Element         `json:"children"`
}

var embedFs embed.FS

var annotationMap map[string][]Element

func InitAnnotation() {
	embedFs = environment.GetEnv().GetEmbed()
	jsonFile, err := embedFs.ReadFile("resources/annotation.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(jsonFile, &annotationMap)
	if err != nil {
		panic(err)
	}
}

func Get(name string) []Element {
	if annotationMap == nil {
		InitAnnotation()
	}
	return annotationMap[name]
}
