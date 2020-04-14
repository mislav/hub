package md2roff

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/russross/blackfriday"
)

// https://github.com/russross/blackfriday/blob/v2/markdown.go
const (
	ParserExtensions = blackfriday.NoIntraEmphasis |
		blackfriday.FencedCode |
		blackfriday.SpaceHeadings |
		blackfriday.AutoHeadingIDs |
		blackfriday.DefinitionLists
)

var (
	backslash     = []byte{'\\'}
	enterVar      = []byte("<var>")
	closeVar      = []byte("</var>")
	tilde         = []byte(`\(ti`)
	htmlEscape    = regexp.MustCompile(`<([A-Za-z][A-Za-z0-9_-]*)>`)
	roffEscape    = regexp.MustCompile(`[&'\_-]`)
	headingEscape = regexp.MustCompile(`["]`)
	titleRe       = regexp.MustCompile(`(?P<name>[A-Za-z][A-Za-z0-9_-]+)\((?P<num>\d)\) -- (?P<title>.+)`)
)

func escape(src []byte, re *regexp.Regexp) []byte {
	return re.ReplaceAllFunc(src, func(c []byte) []byte {
		return append(backslash, c...)
	})
}

func roffText(src []byte) []byte {
	return bytes.Replace(escape(src, roffEscape), []byte{'~'}, tilde, -1)
}

type RoffRenderer struct {
	Manual  string
	Version string
	Date    string
	Title   string
	Name    string
	Section uint8

	listWasTerm bool
}

func (r *RoffRenderer) RenderHeader(buf io.Writer, ast *blackfriday.Node) {
}

func (r *RoffRenderer) RenderFooter(buf io.Writer, ast *blackfriday.Node) {
}

func (r *RoffRenderer) RenderNode(buf io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	if entering {
		switch node.Type {
		case blackfriday.Emph:
			io.WriteString(buf, `\fI`)
		case blackfriday.Strong:
			io.WriteString(buf, `\fB`)
		case blackfriday.Link:
			io.WriteString(buf, `\[la]`)
		case blackfriday.Code:
			io.WriteString(buf, `\fB\fC`)
		case blackfriday.Hardbreak:
			io.WriteString(buf, "\n.br\n")
		case blackfriday.Paragraph:
			if node.Parent.Type != blackfriday.Item {
				io.WriteString(buf, ".P\n")
			} else if node.Parent.FirstChild != node {
				io.WriteString(buf, ".sp\n")
				if node.Prev.Type == blackfriday.List {
					io.WriteString(buf, ".PP\n")
				}
			}
		case blackfriday.CodeBlock:
			io.WriteString(buf, ".PP\n.RS 4\n.nf\n")
		case blackfriday.Item:
			if node.ListFlags&blackfriday.ListTypeDefinition == 0 {
				if node.Parent.ListData.Tight && node.Parent.FirstChild != node {
					io.WriteString(buf, ".sp -1\n")
				}
				if node.Parent.ListData.Tight {
					io.WriteString(buf, ".IP \\(bu 2.3\n")
				} else {
					io.WriteString(buf, ".IP \\(bu 4\n")
				}
			} else {
				if node.ListFlags&blackfriday.ListTypeTerm != 0 {
					io.WriteString(buf, ".PP\n")
				} else {
					io.WriteString(buf, ".RS 4\n")
				}
			}
		case blackfriday.Heading:
			r.renderHeading(buf, node)
			return blackfriday.SkipChildren
		}
	}

	leaf := len(node.Literal) > 0
	if leaf {
		if bytes.Equal(node.Literal, enterVar) {
			io.WriteString(buf, `\fI`)
		} else if bytes.Equal(node.Literal, closeVar) {
			io.WriteString(buf, `\fP`)
		} else {
			buf.Write(roffText(node.Literal))
		}
	}

	if !entering || leaf {
		switch node.Type {
		case blackfriday.Emph,
			blackfriday.Strong:
			io.WriteString(buf, `\fP`)
		case blackfriday.Link:
			io.WriteString(buf, `\[ra]`)
		case blackfriday.Code:
			io.WriteString(buf, `\fR`)
		case blackfriday.CodeBlock:
			io.WriteString(buf, ".fi\n.RE\n")
		case blackfriday.HTMLSpan,
			blackfriday.Del,
			blackfriday.Image:
		case blackfriday.List:
			io.WriteString(buf, ".br\n")
		case blackfriday.Item:
			if node.ListFlags&blackfriday.ListTypeDefinition != 0 &&
				node.ListFlags&blackfriday.ListTypeTerm == 0 {
				io.WriteString(buf, ".RE\n")
			}
		default:
			if !leaf {
				io.WriteString(buf, "\n")
			}
		}
	}

	return blackfriday.GoToNext
}

func textContent(node *blackfriday.Node) []byte {
	var buf bytes.Buffer
	node.Walk(func(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if entering && len(n.Literal) > 0 {
			buf.Write(n.Literal)
		}
		return blackfriday.GoToNext
	})
	return buf.Bytes()
}

func (r *RoffRenderer) renderHeading(buf io.Writer, node *blackfriday.Node) {
	text := textContent(node)
	switch node.HeadingData.Level {
	case 1:
		var name []byte
		var num []byte
		if match := titleRe.FindAllSubmatch(text, 1); match != nil {
			name, num, text = match[0][1], match[0][2], match[0][3]
			r.Name = string(name)
			if sectionNum, err := strconv.Atoi(string(num)); err == nil {
				r.Section = uint8(sectionNum)
			}
			r.Title = string(text)
		}
		fmt.Fprintf(buf, ".TH \"%s\" \"%s\" \"%s\" \"%s\" \"%s\"\n",
			escape(name, headingEscape),
			num,
			escape([]byte(r.Date), headingEscape),
			escape([]byte(r.Version), headingEscape),
			escape([]byte(r.Manual), headingEscape),
		)
		io.WriteString(buf, ".nh\n")   // disable hyphenation
		io.WriteString(buf, ".ad l\n") // disable justification
		io.WriteString(buf, ".SH \"NAME\"\n")
		fmt.Fprintf(buf, "%s \\- %s\n",
			roffText(name),
			roffText(text),
		)
	case 2:
		fmt.Fprintf(buf, ".SH \"%s\"\n", strings.ToUpper(string(escape(text, headingEscape))))
	case 3:
		fmt.Fprintf(buf, ".SS \"%s\"\n", escape(text, headingEscape))
	}
}

func sanitizeInput(src []byte) []byte {
	return htmlEscape.ReplaceAllFunc(src, func(match []byte) []byte {
		res := append(enterVar, match[1:len(match)-1]...)
		return append(res, closeVar...)
	})
}

type renderOption struct {
	renderer blackfriday.Renderer
	buffer   io.Writer
}

func Opt(buffer io.Writer, renderer blackfriday.Renderer) *renderOption {
	return &renderOption{renderer, buffer}
}

func Generate(src []byte, opts ...*renderOption) {
	parser := blackfriday.New(blackfriday.WithExtensions(ParserExtensions))
	ast := parser.Parse(sanitizeInput(src))

	for _, opt := range opts {
		opt.renderer.RenderHeader(opt.buffer, ast)
		ast.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
			return opt.renderer.RenderNode(opt.buffer, node, entering)
		})
		opt.renderer.RenderFooter(opt.buffer, ast)
	}
}
