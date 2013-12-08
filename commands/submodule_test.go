package commands

import (
	"github.com/bmizerany/assert"
	"github.com/jingweno/gh/github"
	"testing"
)

func TestTransformSubmoduleArgs(t *testing.T) {
	github.CreateTestConfig("jingweno", "123")

	args := NewArgs([]string{"submodule", "add", "jingweno/gh", "vendor/gh"})
	transformSubmoduleArgs(args)

	cmds := args.Commands()
	assert.Equal(t, 1, len(cmds))
	assert.Equal(t, "git submodule add git://github.com/jingweno/gh.git vendor/gh", cmds[0].String())

	args = NewArgs([]string{"submodule", "add", "-p", "jingweno/gh",
		"vendor/gh"})
	transformSubmoduleArgs(args)

	cmds = args.Commands()
	assert.Equal(t, 1, len(cmds))
	assert.Equal(t, "git submodule add git@github.com:jingweno/gh.git vendor/gh", cmds[0].String())

	args = NewArgs([]string{"submodule", "add", "-b", "gh", "--name", "gh", "jingweno/gh", "vendor/gh"})
	transformSubmoduleArgs(args)

	cmds = args.Commands()
	assert.Equal(t, 1, len(cmds))
	assert.Equal(t, "git submodule add -b gh --name gh git://github.com/jingweno/gh.git vendor/gh", cmds[0].String())
}
