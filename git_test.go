package main

import (
	"github.com/bmizerany/assert"
	"strings"
	"testing"
)

func TestGitMethods(t *testing.T) {
	gitDir, _ := FetchGitDir()
	assert.T(t, strings.Contains(gitDir, ".git"))

	gitEditor, err := FetchGitEditor()
	if err != nil {
		assert.NotEqual(t, "", gitEditor)
	}

	gitRemote, _ := FetchGitRemote()
	assert.Equal(t, "git@github.com:jingweno/gh.git", gitRemote)

	gitOwner, _ := FetchGitOwner()
	assert.Equal(t, "jingweno", gitOwner)

	gitProject, _ := FetchGitProject()
	assert.Equal(t, "gh", gitProject)

	gitHead, _ := FetchGitHead()
	assert.NotEqual(t, "", gitHead)

	logs, _ := FetchGitCommitLogs("master", "HEAD")
	assert.T(t, len(logs) >= 0)
}
