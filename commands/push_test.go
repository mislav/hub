package commands

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestGetRemotesRef(t *testing.T) {
	args := NewArgs([]string{"push", "origin", "master"})
	remotes, ref := getRemotesRef(args)
	assert.Equal(t, remotes, []string{"origin"})
	assert.Equal(t, ref, "master")

	args = NewArgs([]string{"push", "origin"})
	remotes, ref = getRemotesRef(args)
	assert.Equal(t, remotes, []string{"origin"})
	assert.Equal(t, ref, "")

	args = NewArgs([]string{"push", "origin,experimental", "master"})
	remotes, ref = getRemotesRef(args)
	assert.Equal(t, remotes, []string{"origin", "experimental"})
	assert.Equal(t, ref, "master")

	args = NewArgs([]string{"push", "origin,experimental"})
	remotes, ref = getRemotesRef(args)
	assert.Equal(t, remotes, []string{"origin", "experimental"})
	assert.Equal(t, ref, "")
}
