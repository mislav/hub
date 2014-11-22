package github

import (
	"net/url"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/github/hub/fixtures"
)

func TestGitHubRepo_OriginRemote(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	localRepo, _ := LocalRepo()
	gitRemote, _ := localRepo.OriginRemote()
	assert.Equal(t, "origin", gitRemote.Name)

	u, _ := url.Parse(repo.Remote)
	assert.Equal(t, u, gitRemote.URL)
}

func TestGitHubRepo_remotesForPublish(t *testing.T) {
	url, _ := url.Parse("ssh://git@github.com/Owner/repo.git")
	remotes := []Remote{
		{
			Name: "Owner",
			URL:  url,
		},
	}
	repo := GitHubRepo{remotes}
	remotesForPublish := repo.remotesForPublish("owner")

	assert.Equal(t, 1, len(remotesForPublish))
	assert.Equal(t, "Owner", remotesForPublish[0].Name)
	assert.Equal(t, url.String(), remotesForPublish[0].URL.String())
}
