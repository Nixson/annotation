package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/Nixson/annotation"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"regexp"
	"strings"
)

//go:embed tpl/*.goTpl
var tpls embed.FS

func main() {
	generate(Scan())
}

func Scan() map[string][]annotation.Element {

	fileSystem := os.DirFS(".")
	dirs := make([]string, 0)
	_ = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && path[0:1] != "." {
			dirs = append(dirs, path)
		}
		return nil
	})
	annotations := make([]annotation.Element, 0)
	for _, dir := range dirs {
		d, err := parser.ParseDir(token.NewFileSet(), dir, nil, parser.ParseComments)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for k, f := range d {
			fmt.Println(f, k)
			p := doc.New(f, k, 0)
			for _, tp := range p.Types {
				if tp.Doc != "" {
					annotationE := getAnnotation(tp.Name, tp.Doc)
					if annotationE.Type == "Controller" {
						annotationE.Children = make([]annotation.Element, 0)
						for _, method := range tp.Methods {
							if method.Doc != "" {
								annotationE.Children = append(annotationE.Children, getAnnotation(method.Name, method.Doc))
							}
						}
					}
					annotations = append(annotations, annotationE)
				}
			}
		}

	}
	if len(annotations) > 0 {
		annMap := make(map[string][]annotation.Element)
		annMap["controller"] = get("Controller", annotations)
		annMap["crud"] = get("CRUD", annotations)
		annMap["kafka"] = get("KafkaListen", annotations)
		annotationE, _ := os.Create("resources/annotation.json")
		writr := bufio.NewWriter(annotationE)
		b, err := json.Marshal(annMap)
		if err == nil {
			_, _ = writr.Write(b)
			_ = writr.Flush()
		}
		return annMap
	}
	return nil
}

func get(s string, annotations []annotation.Element) []annotation.Element {
	resp := make([]annotation.Element, 0)
	for _, annotationE := range annotations {
		if annotationE.Type == s {
			resp = append(resp, annotationE)
		}
	}
	return resp
}

func getAnnotation(name string, in string) annotation.Element {
	ann := annotation.Element{
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
		paramsMap[strings.TrimSpace(keyVal[0])] = strings.Trim(strings.TrimSpace(keyVal[1]), `"`)
	}
	return paramsMap
}

var isVar = regexp.MustCompile(`\$\{(.*?)\}`)

func generate(annotationMap map[string][]annotation.Element) {
	controller, ok := annotationMap["controller"]
	if ok {
		//		list := make([]string, 0)
		fileTpl, _ := tpls.ReadFile("controller.goTpl")
		find := isVar.FindStringSubmatch(string(fileTpl))
		fmt.Println(find)
		for _, element := range controller {
			fmt.Println(element)
		}
	}
}
