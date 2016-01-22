package github

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/github/hub/fixtures"
)

func TestGithubRemote_NoPush(t *testing.T) {
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

func TestGithubRemote_GitPlusSsh(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	remoteName := "upstream"
	repo.AddRemote(remoteName, "git+ssh://git@github.com/frozencemetery/python-gssapi", "")

	remotes, err := Remotes()
	assert.Equal(t, nil, err)
	assert.Equal(t, len(remotes), 2)
	assert.Equal(t, remotes[0].Name, remoteName)
	assert.Equal(t, remotes[0].URL.Scheme, "ssh")
	assert.Equal(t, remotes[0].URL.Host, "github.com")
	assert.Equal(t, remotes[0].URL.Path, "/frozencemetery/python-gssapi")
	assert.Equal(t, remotes[1].Name, "origin")
	assert.Equal(t, remotes[1].URL.Path, repo.Remote)
}

func TestGithubRemote_SshPort(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	remoteName := "upstream"
	repo.AddRemote(remoteName, "ssh://git@github.com:22/hakatashi/dotfiles.git", "")

	remotes, err := Remotes()
	assert.Equal(t, nil, err)
	assert.Equal(t, len(remotes), 2)
	assert.Equal(t, remotes[0].Name, remoteName)
	assert.Equal(t, remotes[0].URL.Scheme, "ssh")
	assert.Equal(t, remotes[0].URL.Host, "github.com")
	assert.Equal(t, remotes[0].URL.Path, "/hakatashi/dotfiles.git")
	assert.Equal(t, remotes[1].Name, "origin")
	assert.Equal(t, remotes[1].URL.Path, repo.Remote)
}

func TestGithubRemote_ColonSlash(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	remoteName := "upstream"
	repo.AddRemote(remoteName, "git@github.com:/fatso83/my-project.git", "")

	remotes, err := Remotes()
	assert.Equal(t, nil, err)
	assert.Equal(t, len(remotes), 2)
	assert.Equal(t, remotes[0].Name, remoteName)
	assert.Equal(t, remotes[0].URL.Scheme, "ssh")
	assert.Equal(t, remotes[0].URL.Host, "github.com")
	assert.Equal(t, remotes[0].URL.Path, "/fatso83/my-project.git")
	assert.Equal(t, remotes[1].Name, "origin")
	assert.Equal(t, remotes[1].URL.Path, repo.Remote)
}
