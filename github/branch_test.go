package github

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestBranch_ShortName(t *testing.T) {
	lp, _ := LocalRepo()
	b := Branch{lp, "refs/heads/master"}
	assert.Equal(t, "master", b.ShortName())
}

func TestBranch_LongName(t *testing.T) {
	lp, _ := LocalRepo()

	b := Branch{lp, "refs/heads/master"}
	assert.Equal(t, "heads/master", b.LongName())

	b = Branch{lp, "refs/remotes/origin/master"}
	assert.Equal(t, "origin/master", b.LongName())
}

func TestBranch_RemoteName(t *testing.T) {
	lp, _ := LocalRepo()

	b := Branch{lp, "refs/remotes/origin/master"}
	assert.Equal(t, "origin", b.RemoteName())

	b = Branch{lp, "refs/head/master"}
	assert.Equal(t, "", b.RemoteName())
}

func TestBranch_IsRemote(t *testing.T) {
	lp, _ := LocalRepo()

	b := Branch{lp, "refs/remotes/origin/master"}
	assert.T(t, b.IsRemote())
}
