package channels

import (
	"os"
	"path/filepath"
	"text/template"
)

var tmplDir = filepath.Join(filepath.Dir(os.Args[0]), "templates")

func ParseTemplates(tmplDir string) (*template.Template, error) {
	var tmplFilesPaths []string
	err := filepath.Walk(tmplDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			tmplFilesPaths = append(tmplFilesPaths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return template.ParseFiles(tmplFilesPaths...)
}
