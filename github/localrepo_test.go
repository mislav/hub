package github

import (
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
	assert.Equal(t, repo.Remote, gitRemote.URL.String())
}
