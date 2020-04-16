package github

import (
	"testing"

	"github.com/github/hub/v2/internal/assert"
)

func TestMessageBuilder_multiline_title(t *testing.T) {
	builder := &MessageBuilder{
		Message: `hello
multiline
text

the rest is
description`,
	}

	title, body, err := builder.Extract()
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello multiline text", title)
	assert.Equal(t, "the rest is\ndescription", body)
}
