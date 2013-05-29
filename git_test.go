package main

import (
	"github.com/bmizerany/assert"
	"strings"
	"testing"
)

func TestGitMethods(t *testing.T) {
	git = Git{"git"}
	gitDir, _ := git.Dir()
	assert.T(t, strings.Contains(gitDir, ".git"))

	gitPullReqMsgFile, _ := git.PullReqMsgFile()
	assert.T(t, strings.Contains(gitPullReqMsgFile, "PULLREQ_EDITMSG"))

	gitEditor, err := git.Editor()
	if err == nil {
		assert.NotEqual(t, "", gitEditor)
	}

	gitEditorPath, err := git.EditorPath()
	if err == nil {
		assert.NotEqual(t, "", gitEditorPath)
	}

	gitRemote, _ := git.Remote()
	assert.T(t, strings.Contains(gitRemote, "jingweno/gh.git"))

	gitOwner, _ := git.Owner()
	assert.Equal(t, "jingweno", gitOwner)

	gitProject, _ := git.Project()
	assert.Equal(t, "gh", gitProject)

	gitHead, _ := git.Head()
	assert.NotEqual(t, "", gitHead)

	logs, _ := git.Log("master", "HEAD")
	assert.T(t, len(logs) >= 0)
}
