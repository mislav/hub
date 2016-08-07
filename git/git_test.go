package git

import (
	"fmt"
	"os"
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
	editor := os.Getenv("GIT_EDITOR")
	if err := os.Unsetenv("GIT_EDITOR"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		repo.TearDown()
		if err := os.Setenv("GIT_EDITOR", editor); err != nil {
			t.Fatal(err)
		}
		os.Unsetenv("FOO")
		os.Unsetenv("BAR")
	}()

	os.Setenv("FOO", "hello")
	os.Setenv("BAR", "happy world")

	SetGlobalConfig("core.editor", `$FOO "${BAR}"`)
	gitEditor, err := Editor()
	assert.Equal(t, nil, err)
	assert.Equal(t, `hello "happy world"`, gitEditor)
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

func TestRemotes(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	type remote struct {
		name    string
		url     string
		pushUrl string
	}
	testCases := map[string]remote{
		"testremote1": {
			"testremote1",
			"https://example.com/test1/project1.git",
			"no_push",
		},
		"testremote2": {
			"testremote2",
			"user@example.com:test2/project2.git",
			"http://example.com/project.git",
		},
		"testremote3": {
			"testremote3",
			"https://example.com/test1/project2.git",
			"",
		},
	}

	for _, tc := range testCases {
		repo.AddRemote(tc.name, tc.url, tc.pushUrl)
	}

	remotes, err := Remotes()
	assert.Equal(t, nil, err)

	// In addition to the remotes we added to the repo, repo will
	// also have an additional remote "origin". So add it to the
	// expected cases to test.
	wantCases := map[string]struct{}{
		fmt.Sprintf("origin	%s (fetch)", repo.Remote): {},
		fmt.Sprintf("origin	%s (push)", repo.Remote): {},
		"testremote1	https://example.com/test1/project1.git (fetch)": {},
		"testremote1	no_push (push)": {},
		"testremote2	user@example.com:test2/project2.git (fetch)": {},
		"testremote2	http://example.com/project.git (push)": {},
		"testremote3	https://example.com/test1/project2.git (fetch)": {},
		"testremote3	https://example.com/test1/project2.git (push)": {},
	}

	assert.Equal(t, len(remotes), len(wantCases))
	for _, got := range remotes {
		if _, ok := wantCases[got]; !ok {
			t.Errorf("Unexpected remote: %s", got)
		}
	}
}
