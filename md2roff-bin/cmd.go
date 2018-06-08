package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/github/hub/md2roff"
	flag "github.com/ogier/pflag"
	"github.com/russross/blackfriday"
)

var (
	flagManual,
	flagDate string
)

func init() {
	flag.StringVarP(&flagManual, "manual", "m", "", "MANUAL")
	flag.StringVarP(&flagDate, "date", "d", "", "DATE")
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

	md2roff.Generate(content,
		md2roff.Opt(roffBuf, roff),
		md2roff.Opt(htmlBuf, html),
	)

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
