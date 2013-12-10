package commands

import (
	"bufio"
	"github.com/bmizerany/assert"
	"github.com/jingweno/gh/github"
	"strings"
	"testing"
)

func TestReadTitleAndBody(t *testing.T) {
	message := `A title
A title continues

A body
A body continues
# comment
`
	r := strings.NewReader(message)
	reader := bufio.NewReader(r)
	title, body, err := readTitleAndBodyFrom(reader)
	assert.Equal(t, nil, err)
	assert.Equal(t, "A title A title continues", title)
	assert.Equal(t, "A body\nA body continues", body)
}

func TestParsePullRequestProject(t *testing.T) {
	c := &github.Project{Host: "github.com", Owner: "jingweno", Name: "gh"}

	s := "develop"
	p, ref := parsePullRequestProject(c, s)
	assert.Equal(t, "develop", ref)
	assert.Equal(t, "github.com", p.Host)
	assert.Equal(t, "jingweno", p.Owner)
	assert.Equal(t, "gh", p.Name)

	s = "mojombo:develop"
	p, ref = parsePullRequestProject(c, s)
	assert.Equal(t, "develop", ref)
	assert.Equal(t, "github.com", p.Host)
	assert.Equal(t, "mojombo", p.Owner)
	assert.Equal(t, "gh", p.Name)

	s = "mojombo/jekyll:develop"
	p, ref = parsePullRequestProject(c, s)
	assert.Equal(t, "develop", ref)
	assert.Equal(t, "github.com", p.Host)
	assert.Equal(t, "mojombo", p.Owner)
	assert.Equal(t, "jekyll", p.Name)
}
