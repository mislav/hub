package commands

import (
	"testing"

	"github.com/github/hub/Godeps/_workspace/src/github.com/bmizerany/assert"
)

func TestCommandUseSelf(t *testing.T) {
	c := &Command{Usage: "foo"}

	args := NewArgs([]string{"foo"})

	run, err := lookupCommand(c, args)

	assert.Equal(t, nil, err)
	assert.Equal(t, c, run)
}

func TestCommandUseSubcommand(t *testing.T) {
	c := &Command{Usage: "foo"}
	s := &Command{Usage: "bar"}
	c.Use(s)

	args := NewArgs([]string{"foo", "bar"})

	run, err := lookupCommand(c, args)

	assert.Equal(t, nil, err)
	assert.Equal(t, s, run)
}

func TestCommandUseErrorWhenMissingSubcommand(t *testing.T) {
	c := &Command{Usage: "foo"}
	s := &Command{Usage: "bar"}
	c.Use(s)

	args := NewArgs([]string{"foo", "baz"})

	_, err := lookupCommand(c, args)

	assert.NotEqual(t, nil, err)
}

func TestArgsForCommand(t *testing.T) {
	c := &Command{Usage: "foo"}

	args := NewArgs([]string{"foo", "bar", "baz"})

	lookupCommand(c, args)

	assert.Equal(t, 2, len(args.Params))
}

func TestArgsForSubCommand(t *testing.T) {
	c := &Command{Usage: "foo"}
	s := &Command{Usage: "bar"}
	c.Use(s)

	args := NewArgs([]string{"foo", "bar", "baz"})

	lookupCommand(c, args)

	assert.Equal(t, 1, len(args.Params))
}

func TestFlagsAfterArguments(t *testing.T) {
	c := &Command{Usage: "foo -m MESSAGE ARG1"}

	var flag string
	c.Flag.StringVarP(&flag, "message", "m", "", "MESSAGE")

	args := NewArgs([]string{"foo", "bar", "-m", "baz"})

	c.parseArguments(args)
	assert.Equal(t, "baz", flag)

	assert.Equal(t, 1, len(args.Params))
	assert.Equal(t, "bar", args.LastParam())
}

func TestCommandUsageSubCommands(t *testing.T) {
	f1 := func(c *Command, args *Args) {}
	f2 := func(c *Command, args *Args) {}

	c := &Command{Usage: "foo", Run: f1}
	s := &Command{Key: "bar", Usage: "foo bar", Run: f2}
	c.Use(s)

	usage := c.subCommandsUsage()

	expected := `usage: git foo
   or: git foo bar
`
	assert.Equal(t, expected, usage)
}

func TestCommandUsageSubCommandsPrintOnlyRunnables(t *testing.T) {
	f1 := func(c *Command, args *Args) {}

	c := &Command{Usage: "foo"}
	s := &Command{Key: "bar", Usage: "foo bar", Run: f1}
	c.Use(s)

	usage := c.subCommandsUsage()

	expected := `usage: git foo bar
`
	assert.Equal(t, expected, usage)
}

func TestCommandNameTakeUsage(t *testing.T) {
	c := &Command{Usage: "foo -t -v --foo"}
	assert.Equal(t, "foo", c.Name())
}

func TestCommandNameTakeKey(t *testing.T) {
	c := &Command{Key: "bar", Usage: "foo -t -v --foo"}
	assert.Equal(t, "bar", c.Name())
}

func TestCommandCall(t *testing.T) {
	var result string
	f := func(c *Command, args *Args) { result = args.FirstParam() }

	c := &Command{Usage: "foo", Run: f}
	args := NewArgs([]string{"foo", "bar"})

	c.Call(args)
	assert.Equal(t, "bar", result)
}

func TestCommandHelp(t *testing.T) {
	var result string
	f := func(c *Command, args *Args) { result = args.FirstParam() }
	c := &Command{Usage: "foo", Run: f}
	args := NewArgs([]string{"foo", "-h"})

	c.Call(args)
	assert.Equal(t, "", result)
}

func TestSubCommandCall(t *testing.T) {
	var result string
	f1 := func(c *Command, args *Args) { result = "noop" }
	f2 := func(c *Command, args *Args) { result = args.LastParam() }

	c := &Command{Usage: "foo", Run: f1}
	s := &Command{Key: "bar", Usage: "foo bar", Run: f2}
	c.Use(s)

	args := NewArgs([]string{"foo", "bar", "baz"})

	c.Call(args)
	assert.Equal(t, "baz", result)
}

func TestSubCommandsUsage(t *testing.T) {
	// with subcommand
	f1 := func(c *Command, args *Args) {}
	f2 := func(c *Command, args *Args) {}

	c := &Command{Usage: "foo", Run: f1}
	s := &Command{Key: "bar", Usage: "foo bar", Run: f2}
	c.Use(s)

	usage := c.subCommandsUsage()
	assert.Equal(t, "usage: git foo\n   or: git foo bar\n", usage)

	// no subcommand
	cc := &Command{Usage: "foo", Run: f1}

	usage = cc.subCommandsUsage()
	assert.Equal(t, "usage: git foo\n", usage)
}
