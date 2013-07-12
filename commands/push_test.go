package commands

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestGetRemotesRef(t *testing.T) {
	args := NewArgs([]string{"push", "origin", "master", "--force"})
	remotes, idx := getRemotes(args)
	assert.Equal(t, remotes, []string{"origin"})
	assert.Equal(t, idx, 0)

	args = NewArgs([]string{"push", "origin,experimental", "master", "--force"})
	remotes, idx = getRemotes(args)
	assert.Equal(t, remotes, []string{"origin", "experimental"})
	assert.Equal(t, idx, 0)

}
