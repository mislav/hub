package commands

import (
	"testing"

	"github.com/github/hub/Godeps/_workspace/src/github.com/bmizerany/assert"
)

func TestTransformPrArgs(t *testing.T) {
	var args *Args
	args = NewArgs([]string{"pr", "--apply", "-3", "33"})
	transformPrArgs(args)
	assert.Equal(t, "hub apply -3 https://github.com/github/hub/pull/33", args.ToCmd().String())

	args = NewArgs([]string{"pr", "--am", "-3", "33"})
	transformPrArgs(args)
	assert.Equal(t, "hub am -3 https://github.com/github/hub/pull/33", args.ToCmd().String())

	args = NewArgs([]string{"pr", "--apply", "--not-pr", "33", "-3", "#33"})
	transformPrArgs(args)
	assert.Equal(t, "hub apply 33 -3 https://github.com/github/hub/pull/33", args.ToCmd().String())

	args = NewArgs([]string{"pr", "--browse", "#33", "-a", "Safari"})
	transformPrArgs(args)
	assert.Equal(t, "open https://github.com/github/hub/pull/33 -a Safari", args.ToCmd().String())
}
