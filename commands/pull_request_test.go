package commands

import (
	"testing"

	"github.com/bmizerany/assert"
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

func TestPullRequest_TrimToOneLine(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"A short string.", "A short string."},
		{"On two\nlines.", "On two…"},
		{"With a double\r\nline break.", "With a double…"},
		{"A veeeeeeeeeeeeeeeeeeeeerrrrrrrrrrrrrrrrrrrrrrrrrrrrryyyy looooooooooooong line that is too big to fit.", "A veeeeeeeeeeeeeeeeeeeeerrrrrrrrrrrrrrrrrrrrrrrrrrrrryyyy looooooooooooong line …"},
	}

	for _, test := range tests {
		if got := trimToOneLine(test.input); got != test.expected {
			t.Errorf("trimToOneLine(%q) = %q, want %q", test.input, got, test.expected)
		}
	}
}
