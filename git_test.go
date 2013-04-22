package main

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestGitMethods(t *testing.T) {
	assert.Equal(t, ".git", git.Dir())
	assert.Equal(t, "vim", git.Editor())
	assert.Equal(t, "git@github.com:jingweno/gh.git", git.Remote())
	assert.Equal(t, "jingweno", git.Owner())
	assert.Equal(t, "gh", git.Repo())
	assert.Equal(t, "pull_request", git.CurrentBranch())
}
