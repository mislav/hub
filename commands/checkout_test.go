package commands

import (
	"testing"

	"github.com/github/hub/Godeps/_workspace/src/github.com/bmizerany/assert"
)

func TestReplaceCheckoutParam(t *testing.T) {
	checkoutURL := "https://github.com/github/hub/pull/12"
	args := NewArgs([]string{"checkout", "-b", checkoutURL})
	replaceCheckoutParam(args, checkoutURL, "jingweno", "origin/master")

	cmd := args.ToCmd()
	assert.Equal(t, "git checkout -b --track -B jingweno origin/master", cmd.String())
}
