package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/github/hub/v2/md2roff"
	"github.com/github/hub/v2/utils"
	"github.com/russross/blackfriday"
)

var (
	flagManual,
	flagVersion,
	flagTemplate,
	flagDate string

	xRefRe = regexp.MustCompile(`\b(?P<name>[a-z][\w-]*)\((?P<section>\d)\)`)

	pageIndex map[string]bool
)

func init() {
	pageIndex = make(map[string]bool)
}

type templateData struct {
	Contents string
	Manual   string
	Date     string
	Version  string
	Title    string
	Name     string
	Section  uint8
}

func generateFromFile(mdFile string) error {
	content, err := ioutil.ReadFile(mdFile)
	if err != nil {
		return fmt.Errorf("%s (%q)", err, mdFile)
	}

	roffFile := strings.TrimSuffix(mdFile, ".md")
	roffBuf, err := os.OpenFile(roffFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("%s (%q)", err, roffFile)
	}
	defer roffBuf.Close()

	htmlFile := strings.TrimSuffix(mdFile, ".md") + ".html"
	htmlBuf, err := os.OpenFile(htmlFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("%s (%q)", err, htmlFile)
	}
	defer htmlBuf.Close()

	html := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: blackfriday.HTMLFlagsNone,
	})
	roff := &md2roff.RoffRenderer{
		Manual:  flagManual,
		Version: flagVersion,
		Date:    flagDate,
	}

	htmlGenBuf := &bytes.Buffer{}
	var htmlWriteBuf io.Writer = htmlBuf
	if flagTemplate != "" {
		htmlWriteBuf = htmlGenBuf
	}

	md2roff.Generate(content,
		md2roff.Opt(roffBuf, roff),
		md2roff.Opt(htmlWriteBuf, html),
	)

	if flagTemplate != "" {
		htmlGenBytes, err := ioutil.ReadAll(htmlGenBuf)
		if err != nil {
			return fmt.Errorf("%s [%s]", err, "htmlGenBuf")
		}
		content := ""
		if contentLines := strings.Split(string(htmlGenBytes), "\n"); len(contentLines) > 1 {
			content = strings.Join(contentLines[1:], "\n")
		}

		currentPage := fmt.Sprintf("%s(%d)", roff.Name, roff.Section)
		content = xRefRe.ReplaceAllStringFunc(content, func(match string) string {
			if match == currentPage {
				return match
			}
			matches := xRefRe.FindAllStringSubmatch(match, 1)
			fileName := fmt.Sprintf("%s.%s", matches[0][1], matches[0][2])
			if pageIndex[fileName] {
				return fmt.Sprintf(`<a href="./%s.html">%s</a>`, fileName, match)
			}
			return match
		})

		tmplData := templateData{
			Manual:   flagManual,
			Date:     flagDate,
			Contents: content,
			Title:    roff.Title,
			Section:  roff.Section,
			Name:     roff.Name,
			Version:  flagVersion,
		}

		templateFile, err := ioutil.ReadFile(flagTemplate)
		if err != nil {
			return fmt.Errorf("%s (%q)", err, flagTemplate)
		}
		tmpl, err := template.New("test").Parse(string(templateFile))
		if err != nil {
			return err
		}
		err = tmpl.Execute(htmlBuf, tmplData)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	p := utils.NewArgsParserWithUsage(`
		--manual NAME
		--version STR
		--template FILE
		--date DATE
	`)
	files, err := p.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}
	flagManual = p.Value("--manual")
	flagVersion = p.Value("--version")
	flagTemplate = p.Value("--template")
	flagDate = p.Value("--date")

	for _, infile := range files {
		name := path.Base(infile)
		name = strings.TrimSuffix(name, ".md")
		pageIndex[name] = true
	}

	for _, infile := range files {
		err := generateFromFile(infile)
		if err != nil {
			panic(err)
		}
	}
}
