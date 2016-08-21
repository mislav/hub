package commands

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestRenderReleaseTpl(t *testing.T) {
	msg, err := renderReleaseTpl("Creating", "#", "1.0", "github/hub", "master")
	assert.Equal(t, nil, err)

	expMsg := `
# Creating release 1.0 for github/hub from master
#
# Write a message for this release. The first block of
# text is the title and the rest is the description.`
	assert.Equal(t, expMsg, msg)
}
