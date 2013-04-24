package main

import (
	"github.com/bmizerany/assert"
	"strings"
	"testing"
)

func TestGitMethods(t *testing.T) {
	assert.T(t, strings.Contains(git.Dir(), ".git"))
	assert.Equal(t, "vim", git.Editor())
	assert.Equal(t, "git@github.com:jingweno/gh.git", git.Remote())
	assert.Equal(t, "jingweno", git.Owner())
	assert.Equal(t, "gh", git.Repo())
	assert.Equal(t, "pull_request", git.CurrentBranch())
}
