package github

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/github/hub/fixtures"
)

func TestOriginRemote(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	gitRemote, _ := OriginRemote()
	assert.Equal(t, "origin", gitRemote.Name)
	assert.Equal(t, repo.Remote, gitRemote.URL.String())
}
