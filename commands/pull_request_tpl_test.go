package commands

import (
	"testing"

	"github.com/github/hub/Godeps/_workspace/src/github.com/bmizerany/assert"
)

func TestRenderPullRequestTpl(t *testing.T) {
	msg, err := renderPullRequestTpl("init", "#", "base", "head", "one\ntwo")
	assert.Equal(t, nil, err)

	expMsg := `init

# Requesting a pull to base from head
#
# Write a message for this pull request. The first block
# of text is the title and the rest is description.
#
# Changes:
#
# one
# two`
	assert.Equal(t, expMsg, msg)
}
