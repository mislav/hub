package git

import (
	"os"
	"strings"
	"testing"

	"github.com/github/hub/Godeps/_workspace/src/github.com/bmizerany/assert"
	"github.com/github/hub/fixtures"
)

func TestGitDir(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	gitDir, _ := Dir()
	assert.T(t, strings.Contains(gitDir, ".git"))
}

func TestGitEditor(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	SetGlobalConfig("core.editor", "foo")
	gitEditor, err := Editor()
	assert.Equal(t, nil, err)
	assert.Equal(t, "foo", gitEditor)
}

func TestGitLog(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	log, err := Log("08f4b7b6513dffc6245857e497cfd6101dc47818", "9b5a719a3d76ac9dc2fa635d9b1f34fd73994c06")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", log)
}

func TestGitRef(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	ref := "08f4b7b6513dffc6245857e497cfd6101dc47818"
	gitRef, err := Ref(ref)
	assert.Equal(t, nil, err)
	assert.Equal(t, ref, gitRef)
}

func TestGitRefList(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	refList, err := RefList("08f4b7b6513dffc6245857e497cfd6101dc47818", "9b5a719a3d76ac9dc2fa635d9b1f34fd73994c06")
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(refList))

	assert.Equal(t, "9b5a719a3d76ac9dc2fa635d9b1f34fd73994c06", refList[0])
}

func TestGitShow(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	output, err := Show("9b5a719a3d76ac9dc2fa635d9b1f34fd73994c06")
	assert.Equal(t, nil, err)
	assert.Equal(t, "First comment\n\nMore comment", output)
}

func TestGitConfig(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	v, err := GlobalConfig("gh.test")
	assert.NotEqual(t, nil, err)

	SetGlobalConfig("gh.test", "1")
	v, err = GlobalConfig("gh.test")
	assert.Equal(t, nil, err)
	assert.Equal(t, "1", v)
}

func TestGitSignatureWithConfig(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	err := SetConfig("user.name", "Some Hacker")
	assert.Equal(t, nil, err)
	err = SetConfig("user.email", "hacker@example.com")
	assert.Equal(t, nil, err)

	s, err := AuthorSignature()
	assert.Equal(t, nil, err)
	assert.Equal(t, "Signed-off-by: Some Hacker <hacker@example.com>", s)
}

func TestGitSignatureWithEnv(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	os.Setenv("GIT_AUTHOR_NAME", "Some Hacker")
	defer os.Setenv("GIT_AUTHOR_NAME", "")
	os.Setenv("GIT_AUTHOR_EMAIL", "hacker@example.com")
	defer os.Setenv("GIT_AUTHOR_EMAIL", "")

	s, err := AuthorSignature()
	assert.Equal(t, nil, err)
	assert.Equal(t, "Signed-off-by: Some Hacker <hacker@example.com>", s)
}
