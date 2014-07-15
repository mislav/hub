package git

import (
	"strings"
	"testing"

	"github.com/bmizerany/assert"
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

func TestGitLog2(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	log, err := Log2("08f4b7b6513dffc6245857e497cfd6101dc47818", "9b5a719a3d76ac9dc2fa635d9b1f34fd73994c06")
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, log)
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
