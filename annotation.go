package annotation

import (
	"embed"
	"encoding/json"
)

type Annotation struct{}
type Element struct {
	Type       string            `json:"type"`
	StructName string            `json:"structName"`
	Parameters map[string]string `json:"parameters"`
	Children   []Element         `json:"children"`
}

var annotation Annotation
var embedFs embed.FS

var annotationMap map[string][]Element

func InitAnnotation(emb embed.FS) *Annotation {
	embedFs = emb
	jsonFile, err := embedFs.ReadFile("resources/annotation.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(jsonFile, &annotationMap)
	if err != nil {
		panic(err)
	}
	annotation = Annotation{}
	return &annotation
}

func (a *Annotation) Get(name string) []Element {
	return annotationMap[name]
}
