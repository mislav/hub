package commands

import (
	"flag"
	"fmt"
	"strings"
)

type Command struct {
	Run  func(cmd *Command, args []string)
	Flag flag.FlagSet

	Usage string
	Short string
	Long  string
}

func (c *Command) PrintUsage() {
	if c.Runnable() {
		fmt.Printf("Usage: gh %s\n\n", c.Usage)
	}
	fmt.Println(strings.Trim(c.Long, "\n"))
}

func (c *Command) Name() string {
	name := c.Usage
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}

func (c *Command) List() bool {
	return c.Short != ""
}

var All = []*Command{
	cmdPull,
	cmdFork,
	cmdCi,
	cmdBrowse,
	cmdCompare,
	cmdHelp,
	cmdVersion,
}
