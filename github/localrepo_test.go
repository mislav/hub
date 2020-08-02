package github

import (
	"net/url"
	"testing"

	"github.com/github/hub/v2/internal/assert"
)

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
