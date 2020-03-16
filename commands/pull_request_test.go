package commands

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/github/hub/fixtures"
	"github.com/github/hub/github"
)

func TestPullRequest_ParsePullRequestProject(t *testing.T) {
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

func TestPrepareMessageFromCommit(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	prTemplate := "My PR Template Rocks."
	repo.AddFile("test.git/.github/pull_request_template.md", prTemplate)

	message, err := prepareMessageFromCommit("9b5a719a3d76ac9dc2fa635d9b1f34fd73994c06")
	assert.Equal(t, nil, err)
	assert.Equal(t, "First comment\n\nMore comment"+"\n\n\n"+prTemplate, message)
}
