package github

import (
	"testing"

	"github.com/github/hub/Godeps/_workspace/src/github.com/bmizerany/assert"
	"github.com/github/hub/fixtures"
)

func TestGithubRepo_Remotes(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	remoteName := "upstream"
	repo.AddRemote(remoteName, "user@example.com:test/project.git", "no_push")

	remotes, err := Remotes()
	assert.Equal(t, nil, err)
	assert.Equal(t, len(remotes), 2)
	assert.Equal(t, remotes[0].Name, remoteName)
	assert.Equal(t, remotes[0].URL.Scheme, "ssh")
	assert.Equal(t, remotes[0].URL.Host, "example.com")
	assert.Equal(t, remotes[0].URL.Path, "/test/project.git")
	assert.Equal(t, remotes[1].Name, "origin")
	assert.Equal(t, remotes[1].URL.Path, repo.Remote)
}
