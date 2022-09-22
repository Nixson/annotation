package annotation

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"strings"
)

type Element struct {
	Type       string            `json:"type"`
	StructName string            `json:"structName"`
	Parameters map[string]string `json:"parameters"`
	Children   []Element         `json:"children"`
}

func Scan() {

	fileSystem := os.DirFS(".")
	dirs := make([]string, 0)
	_ = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && path[0:1] != "." {
			dirs = append(dirs, path)
		}
		return nil
	})
	annotations := make([]Element, 0)
	for _, dir := range dirs {
		d, err := parser.ParseDir(token.NewFileSet(), dir, nil, parser.ParseComments)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for k, f := range d {
			p := doc.New(f, k, 0)
			for _, tp := range p.Types {
				if tp.Doc != "" {
					annotation := getAnnotation(tp.Name, tp.Doc)
					if annotation.Type == "Controller" {
						annotation.Children = make([]Element, 0)
						for _, method := range tp.Methods {
							if method.Doc != "" {
								annotation.Children = append(annotation.Children, getAnnotation(method.Name, method.Doc))
							}
						}
					}
					annotations = append(annotations, annotation)
				}
			}
		}

	}
	if len(annotations) > 0 {
		annMap := make(map[string][]Element)
		annMap["controller"] = get("Controller", annotations)
		annMap["crud"] = get("CRUD", annotations)
		annotation, _ := os.Create("resources/annotation.json")
		writr := bufio.NewWriter(annotation)
		b, err := json.Marshal(annMap)
		if err == nil {
			_, _ = writr.Write(b)
			_ = writr.Flush()
		}
	}
}

func get(s string, annotations []Element) []Element {
	resp := make([]Element, 0)
	for _, annotation := range annotations {
		if annotation.Type == s {
			resp = append(resp, annotation)
		}
	}
	return resp
}

func getAnnotation(name string, in string) Element {
	ann := Element{
		StructName: name,
	}
	sep := strings.Split(in, "\n")
	for _, str := range sep {
		if strings.Contains(str, "@") {
			titleApp := strings.Split(str, "@")
			title := strings.TrimSpace(titleApp[1])
			if strings.Contains(title, "(") {
				strName := strings.Split(title, "(")
				ann.Type = strName[0]
				ann.Parameters = parseParams(strName[1])
			} else {
				ann.Type = title
			}
		}
	}
	return ann
}

func parseParams(s string) map[string]string {
	paramsMap := make(map[string]string)
	s = s[:len(s)-1]
	sep := strings.Split(s, ",")
	for _, substr := range sep {
		substr = strings.TrimSpace(substr)
		if !strings.Contains(substr, "=") {
			continue
		}
		keyVal := strings.Split(substr, "=")
		paramsMap[strings.TrimSpace(keyVal[0])] = strings.TrimSpace(keyVal[1])
	}
	return paramsMap
}
