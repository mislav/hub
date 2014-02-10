package github

import (
	"os"
	"testing"

	"github.com/bmizerany/assert"
)

type TestRepo struct {
	Pwd string
}

func (g *TestRepo) Setup() {
	g.Pwd, _ = os.Getwd()
	os.Chdir("../fixtures/test.git")
}

func (g *TestRepo) TearDown() {
	os.Chdir(g.Pwd)
}

func setupRepo() *TestRepo {
	repo := &TestRepo{}
	repo.Setup()

	return repo
}

func TestOriginRemote(t *testing.T) {
	repo := setupRepo()
	defer repo.TearDown()

	gitRemote, _ := OriginRemote()
	assert.Equal(t, "origin", gitRemote.Name)
	assert.Equal(t, "https://github.com/test/test.git.git", gitRemote.URL.String())
}
