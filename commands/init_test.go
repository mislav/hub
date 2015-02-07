package commands

import (
	"os"
	"regexp"
	"testing"

	"github.com/github/hub/Godeps/_workspace/src/github.com/bmizerany/assert"
	"github.com/github/hub/github"
)

func TestTransformInitArgs(t *testing.T) {
	os.Setenv("HUB_PROTOCOL", "git")
	github.CreateTestConfigs("jingweno", "123")

	args := NewArgs([]string{"init"})
	err := transformInitArgs(args)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, args.IsParamsEmpty())

	args = NewArgs([]string{"init", "foo"})
	err = transformInitArgs(args)

	assert.Equal(t, nil, err)
	commands := args.Commands()
	assert.Equal(t, 1, len(commands))
	assert.Equal(t, "git init foo", commands[0].String())

	args = NewArgs([]string{"init", "-g", "--bare"})
	err = transformInitArgs(args)

	assert.Equal(t, nil, err)
	commands = args.Commands()
	assert.Equal(t, 2, len(commands))
	assert.Equal(t, "git init --bare", commands[0].String())
	reg := regexp.MustCompile("git remote add origin git@github.com:jingweno/.+\\.git")
	assert.T(t, reg.MatchString(commands[1].String()))
}
