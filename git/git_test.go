package git

import (
	"os"
	"strings"
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

func TestGitDir(t *testing.T) {
	repo := setupRepo()
	defer repo.TearDown()

	gitDir, _ := Dir()
	assert.T(t, strings.Contains(gitDir, ".git"))
}

func TestGitEditor(t *testing.T) {
	gitEditor, err := Editor()
	if err == nil {
		assert.NotEqual(t, "", gitEditor)
	}
}

func TestGitLog(t *testing.T) {
	repo := setupRepo()
	defer repo.TearDown()

	log, err := Log("1dbff497d642562805323c5c2cccd4adc4a83b36", "5196494806847d5233d877517a79b6ce8b33f5f7")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", log)
}

func TestGitRef(t *testing.T) {
	repo := setupRepo()
	defer repo.TearDown()

	ref := "1dbff497d642562805323c5c2cccd4adc4a83b36"
	gitRef, err := Ref(ref)
	assert.Equal(t, nil, err)
	assert.Equal(t, ref, gitRef)
}

func TestGitRefList(t *testing.T) {
	repo := setupRepo()
	defer repo.TearDown()

	refList, err := RefList("1dbff497d642562805323c5c2cccd4adc4a83b36", "5196494806847d5233d877517a79b6ce8b33f5f7")
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(refList))

	assert.Equal(t, "5196494806847d5233d877517a79b6ce8b33f5f7", refList[0])
}

func TestGitShow(t *testing.T) {
	repo := setupRepo()
	defer repo.TearDown()

	output, err := Show("5196494806847d5233d877517a79b6ce8b33f5f7")
	assert.Equal(t, nil, err)
	assert.Equal(t, "Test comment\n\nMore comment", output)
}

func TestGitConfig(t *testing.T) {
	defer UnsetGlobalConfig("gh.test")

	v, err := GlobalConfig("gh.test")
	assert.NotEqual(t, nil, err)

	SetGlobalConfig("gh.test", "1")
	v, err = GlobalConfig("gh.test")
	assert.Equal(t, nil, err)
	assert.Equal(t, "1", v)
}
