package main

import (
	"github.com/bmizerany/assert"
	"strings"
	"testing"
)

func TestGitMethods(t *testing.T) {
	assert.T(t, strings.Contains(FetchGitDir(), ".git"))
	assert.Equal(t, "vim", FetchGitEditor())
	assert.Equal(t, "git@github.com:jingweno/gh.git", FetchGitRemote())
	assert.Equal(t, "jingweno", FetchGitOwner())
	assert.Equal(t, "gh", FetchGitProject())
	assert.Equal(t, "pull_request", FetchGitHead())
	logs := FetchGitCommitLogs("master", "HEAD")
	assert.T(t, len(logs) > 0)
}
