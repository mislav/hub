package main

import (
	"github.com/bmizerany/assert"
	"strings"
	"testing"
)

func setupGit() *Git {
	return &Git{"git"}
}

func TestGitDir(t *testing.T) {
	git := setupGit()
	gitDir, _ := git.Dir()
	assert.T(t, strings.Contains(gitDir, ".git"))
}

func TestGitPullReqMsgFile(t *testing.T) {
	git := setupGit()
	gitPullReqMsgFile, _ := git.PullReqMsgFile()
	assert.T(t, strings.Contains(gitPullReqMsgFile, "PULLREQ_EDITMSG"))
}

func TestGitEditor(t *testing.T) {
	git := setupGit()
	gitEditor, err := git.Editor()
	if err == nil {
		assert.NotEqual(t, "", gitEditor)
	}
}

func TestGitEditorPath(t *testing.T) {
	git := setupGit()
	gitEditorPath, err := git.EditorPath()
	if err == nil {
		assert.NotEqual(t, "", gitEditorPath)
	}
}

func TestGitRemote(t *testing.T) {
	git := setupGit()
	gitRemote, _ := git.Remote()
	assert.T(t, strings.Contains(gitRemote, "jingweno/gh.git"))
}

func TestGitHead(t *testing.T) {
	git := setupGit()
	gitHead, _ := git.Head()
	assert.NotEqual(t, "", gitHead)
}

func TestGitLog(t *testing.T) {
	git := setupGit()
	logs, _ := git.Log("master", "HEAD")
	assert.T(t, len(logs) >= 0)
}

func TestGitRef(t *testing.T) {
	git := setupGit()
	gitRef, _ := git.Ref("HEAD")
	assert.NotEqual(t, "", gitRef)
}
