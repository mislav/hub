package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/github/hub/md2roff"
	flag "github.com/ogier/pflag"
	"github.com/russross/blackfriday"
)

var (
	flagManual,
	flagVersion,
	flagTemplate,
	flagDate string

	pageIndex map[string]bool
)

func init() {
	flag.StringVarP(&flagManual, "manual", "m", "", "MANUAL")
	flag.StringVarP(&flagVersion, "version", "", "", "VERSION")
	flag.StringVarP(&flagTemplate, "template", "t", "", "TEMPLATE")
	flag.StringVarP(&flagDate, "date", "d", "", "DATE")
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
		return err
	}

	roffFile := strings.TrimSuffix(mdFile, ".md")
	roffBuf, err := os.OpenFile(roffFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer roffBuf.Close()

	htmlFile := strings.TrimSuffix(mdFile, ".md") + ".html"
	htmlBuf, err := os.OpenFile(htmlFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer htmlBuf.Close()

	html := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: blackfriday.HTMLFlagsNone,
	})
	roff := &md2roff.RoffRenderer{
		Manual: flagManual,
		Date:   flagDate,
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
			return err
		}
		content := ""
		if contentLines := strings.Split(string(htmlGenBytes), "\n"); len(contentLines) > 1 {
			content = strings.Join(contentLines[1:], "\n")
		}

		tmplData := templateData{
			Manual:   roff.Manual,
			Date:     roff.Date,
			Contents: content,
			Title:    roff.Title,
			Section:  roff.Section,
			Name:     roff.Name,
			Version:  flagVersion,
		}

		templateFile, err := ioutil.ReadFile(flagTemplate)
		if err != nil {
			return err
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
	flag.Parse()

	for _, infile := range flag.Args() {
		err := generateFromFile(infile)
		if err != nil {
			panic(err)
		}
	}
}
