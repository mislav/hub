package github

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestBranch_ShortName(t *testing.T) {
	b := Branch{"refs/heads/master"}
	assert.Equal(t, "master", b.ShortName())
}

func TestBranch_LongName(t *testing.T) {
	b := Branch{"refs/heads/master"}
	assert.Equal(t, "heads/master", b.LongName())

	b = Branch{"refs/remotes/origin/master"}
	assert.Equal(t, "origin/master", b.LongName())
}

func TestBranch_RemoveName(t *testing.T) {
	b := Branch{"refs/remotes/origin/master"}
	assert.Equal(t, "origin", b.RemoteName())
}

func TestBranch_IsRemote(t *testing.T) {
	b := Branch{"refs/remotes/origin/master"}
	assert.T(t, b.IsRemote())
}
