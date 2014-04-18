package git

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestURL_ParseURL(t *testing.T) {
	u, err := ParseURL("https://github.com/octokit/go-octokit.git")
	assert.Equal(t, nil, err)
	assert.Equal(t, "github.com", u.Host)
	assert.Equal(t, "https", u.Scheme)
	assert.Equal(t, "/octokit/go-octokit.git", u.Path)

	u, err = ParseURL("git://github.com/octokit/go-octokit.git")
	assert.Equal(t, nil, err)
	assert.Equal(t, "github.com", u.Host)
	assert.Equal(t, "git", u.Scheme)
	assert.Equal(t, "/octokit/go-octokit.git", u.Path)

	u, err = ParseURL("git@github.com:lostisland/go-sawyer.git")
	assert.Equal(t, nil, err)
	assert.Equal(t, "github.com", u.Host)
	assert.Equal(t, "ssh", u.Scheme)
	assert.Equal(t, "git", u.User.Username())
	assert.Equal(t, "/lostisland/go-sawyer.git", u.Path)
}
