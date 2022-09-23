package annotation

import (
	"embed"
	"encoding/json"
	"github.com/Nixson/environment"
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

func InitAnnotation() *Annotation {
	embedFs = environment.GetEnv().GetEmbed()
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
