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
	editor := os.Getenv("GIT_EDITOR")
	if err := os.Unsetenv("GIT_EDITOR"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		repo.TearDown()
		if err := os.Setenv("GIT_EDITOR", editor); err != nil {
			t.Fatal(err)
		}
	}()

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

	v, err := GlobalConfig("hub.test")
	assert.NotEqual(t, nil, err)

	SetGlobalConfig("hub.test", "1")
	v, err = GlobalConfig("hub.test")
	assert.Equal(t, nil, err)
	assert.Equal(t, "1", v)

	SetGlobalConfig("hub.test", "")
	v, err = GlobalConfig("hub.test")
	assert.Equal(t, nil, err)
	assert.Equal(t, "", v)
}
