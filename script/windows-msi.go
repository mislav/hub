package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"text/template"
)

var nonWordRE = regexp.MustCompile(`\W+`)

var funcs = template.FuncMap{
	"glob": func(p string) ([]string, error) {
		return filepath.Glob(p)
	},
	"id": func(p string) string {
		return nonWordRE.ReplaceAllString(p, "_")
	},
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("needs version argument")
	}
	version := os.Args[1]

	inFile := "script/windows-msi.wxs.tmpl"
	t, err := template.New(path.Base(inFile)).Funcs(funcs).ParseFiles(inFile)
	if err != nil {
		log.Fatal(err)
	}

	data := &struct {
		Version string
	}{
		Version: version,
	}

	if err := t.Execute(os.Stdout, data); err != nil {
		log.Fatal(err)
	}
}
