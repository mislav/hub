package commands

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/github/hub/ui"
)

func TestMain(m *testing.M) {
	ui.Default = ui.Console{Stdout: ioutil.Discard, Stderr: ioutil.Discard}
	os.Exit(m.Run())
}

func TestCommandUseSelf(t *testing.T) {
	c := &Command{Usage: "foo"}

	args := NewArgs([]string{"foo"})

	run, err := c.lookupSubCommand(args)

	assert.Equal(t, nil, err)
	assert.Equal(t, c, run)
}

func TestCommandUseSubcommand(t *testing.T) {
	c := &Command{Usage: "foo"}
	s := &Command{Usage: "bar"}
	c.Use(s)

	args := NewArgs([]string{"foo", "bar"})

	run, err := c.lookupSubCommand(args)

	assert.Equal(t, nil, err)
	assert.Equal(t, s, run)
}

func TestCommandUseErrorWhenMissingSubcommand(t *testing.T) {
	c := &Command{Usage: "foo"}
	s := &Command{Usage: "bar"}
	c.Use(s)

	args := NewArgs([]string{"foo", "baz"})

	_, err := c.lookupSubCommand(args)

	assert.NotEqual(t, nil, err)
}

func TestArgsForCommand(t *testing.T) {
	c := &Command{Usage: "foo"}

	args := NewArgs([]string{"foo", "bar", "baz"})

	c.lookupSubCommand(args)

	assert.Equal(t, 2, len(args.Params))
}

func TestArgsForSubCommand(t *testing.T) {
	c := &Command{Usage: "foo"}
	s := &Command{Usage: "bar"}
	c.Use(s)

	args := NewArgs([]string{"foo", "bar", "baz"})

	c.lookupSubCommand(args)

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
