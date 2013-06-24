package git

import (
	"github.com/bmizerany/assert"
	"strings"
	"testing"
)

func TestGitDir(t *testing.T) {
	gitDir, _ := Dir()
	assert.T(t, strings.Contains(gitDir, ".git"))
}

func TestGitPullReqMsgFile(t *testing.T) {
	gitPullReqMsgFile, _ := PullReqMsgFile()
	assert.T(t, strings.Contains(gitPullReqMsgFile, "PULLREQ_EDITMSG"))
}

func TestGitEditor(t *testing.T) {
	gitEditor, err := Editor()
	if err == nil {
		assert.NotEqual(t, "", gitEditor)
	}
}

func TestGitEditorPath(t *testing.T) {
	gitEditorPath, err := EditorPath()
	if err == nil {
		assert.NotEqual(t, "", gitEditorPath)
	}
}

func TestGitRemote(t *testing.T) {
	gitRemote, _ := OriginRemote()
	assert.Equal(t, "origin", gitRemote.Name)
	assert.T(t, strings.Contains(gitRemote.URL, "gh"))
}

func TestGitHead(t *testing.T) {
	gitHead, _ := Head()
	assert.NotEqual(t, "", gitHead)
}

func TestGitLog(t *testing.T) {
	logs, _ := Log("master", "HEAD")
	assert.T(t, len(logs) >= 0)
}

func TestGitRef(t *testing.T) {
	gitRef, err := Ref("master")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", gitRef)
}
