package commands

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/github/hub/cmd"
)

func TestRunner_splitAliasCmd(t *testing.T) {
	words, err := splitAliasCmd("!source ~/.zshrc")
	assert.NotEqual(t, nil, err)

	words, err = splitAliasCmd("log --pretty=oneline --abbrev-commit --graph --decorate")
	assert.Equal(t, nil, err)
	assert.Equal(t, 5, len(words))

	words, err = splitAliasCmd("")
	assert.NotEqual(t, nil, err)
}

func TestRunnerUseCommands(t *testing.T) {
	r := &Runner{
		commands: make(map[string]*Command),
		execute:  func([]*cmd.Cmd, bool) error { return nil },
	}
	c := &Command{Usage: "foo"}
	r.Use(c)

	assert.Equal(t, c, r.Lookup("foo"))
}

func TestRunnerCallCommands(t *testing.T) {
	var result string
	f := func(c *Command, args *Args) {
		result = args.FirstParam()
		args.Replace("git", "version", "")
	}

	r := &Runner{
		commands: make(map[string]*Command),
		execute:  func([]*cmd.Cmd, bool) error { return nil },
	}
	c := &Command{Usage: "foo", Run: f}
	r.Use(c)

	args := NewArgs([]string{"foo", "bar"})
	err := r.Call(c, args)

	assert.Equal(t, 0, err.ExitCode)
	assert.Equal(t, "bar", result)
}
