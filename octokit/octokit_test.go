package octokit

import (
	"bytes"
	"github.com/bmizerany/assert"
	"testing"
)

func TestGet(t *testing.T) {
	c := NewClient()
	body, err := c.get("repos/jingweno/gh/commits", nil)

	assert.Equal(t, nil, err)
	assert.T(t, len(body) > 0)
}

func TestPost(t *testing.T) {
	content := "# title"
	c := NewClient()

	body, err := c.post("markdown/raw", nil, bytes.NewBufferString(content))
	assert.Equal(t, "Invalid request media type (expecting 'text/plain')", err.Error())

	headers := make(map[string]string)
	headers["Content-Type"] = "text/plain"
	body, err = c.post("markdown/raw", headers, bytes.NewBufferString(content))

	assert.Equal(t, nil, err)
	expectBody := "<h1>\n<a name=\"title\" class=\"anchor\" href=\"#title\"><span class=\"octicon octicon-link\"></span></a>title</h1>"
	assert.Equal(t, expectBody, string(body))
}
