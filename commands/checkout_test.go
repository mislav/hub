package commands

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestReplaceCheckoutParam(t *testing.T) {
	checkoutURL := "https://github.com/github/hub/pull/12"
	args := NewArgs([]string{"checkout", checkoutURL})
	replaceCheckoutParam(args, checkoutURL, "jingweno", "origin/master")

	cmd := args.ToCmd()
	assert.Equal(t, "git checkout --track -B jingweno origin/master", cmd.String())
}

func TestTransformCheckoutArgs(t *testing.T) {
	args := NewArgs([]string{"checkout", "-b", "https://github.com/github/hub/pull/12"})
	err := transformCheckoutArgs(args)

	assert.Equal(t, "Unsupported flag -b when checking out pull request", err.Error())

	args = NewArgs([]string{"checkout", "--orphan", "https://github.com/github/hub/pull/12"})
	err = transformCheckoutArgs(args)

	assert.Equal(t, "Unsupported flag --orphan when checking out pull request", err.Error())
}
