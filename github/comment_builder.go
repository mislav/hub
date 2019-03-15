package github

import (
	"regexp"
)

type CommentBuilder struct {
	Filename string
	Message  string
	Edit     bool
	editor   *Editor
}

func (b *CommentBuilder) Extract() (body string, err error) {
	body = b.Message

	if b.Edit {
		b.editor, err = NewEditor(b.Filename, "", body)
		if err != nil {
			return
		}
		body, err = b.editor.EditContent()
		if err != nil {
			return
		}
	} else {
		nl := regexp.MustCompile(`\r?\n`)
		body = nl.ReplaceAllString(body, "\n")
	}

	if body == "" {
		defer b.Cleanup()
	}

	return
}

func (b *CommentBuilder) Cleanup() {
	if b.editor != nil {
		b.editor.DeleteFile()
	}
}
