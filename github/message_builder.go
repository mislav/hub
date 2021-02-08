package github

import (
	"regexp"
	"strings"
)

type MessageBuilder struct {
	Title             string
	Filename          string
	Message           string
	Edit              bool
	commentedSections []string
	editor            *Editor
}

func (b *MessageBuilder) AddCommentedSection(section string) {
	b.commentedSections = append(b.commentedSections, section)
}

func (b *MessageBuilder) Extract() (title, body string, err error) {
	content := b.Message

	if b.Edit {
		b.editor, err = NewEditor(b.Filename, b.Title, content)
		if err != nil {
			return
		}
		for _, section := range b.commentedSections {
			b.editor.AddCommentedSection(section)
		}
		content, err = b.editor.EditContent()
		if err != nil {
			return
		}
	} else {
		nl := regexp.MustCompile(`\r?\n`)
		content = nl.ReplaceAllString(content, "\n")
	}

	title, body = SplitTitleBody(content)
	if title == "" {
		defer b.Cleanup()
	}

	return
}

func (b *MessageBuilder) Cleanup() {
	if b.editor != nil {
		b.editor.DeleteFile()
	}
}

func SplitTitleBody(content string) (title string, body string) {
	parts := strings.SplitN(content, "\n\n", 2)
	if len(parts) >= 1 {
		title = strings.TrimSpace(strings.Replace(parts[0], "\n", " ", -1))
	}
	if len(parts) >= 2 {
		body = strings.TrimSpace(parts[1])
	}
	return
}
