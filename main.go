package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Command struct {
	Run  func(cmd *Command, args []string)
	Flag flag.FlagSet

	Usage string
	Short string
	Long  string
}

func (c *Command) printUsage() {
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

var commands = []*Command{
	cmdPullRequest,
	cmdHelp,
}

var gh = NewGitHub(os.Getenv("HOME") + "/.config/gh")

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		usage()
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() {
				cmd.printUsage()
			}
			if err := cmd.Flag.Parse(args[1:]); err != nil {
				os.Exit(2)
			}
			cmd.Run(cmd, cmd.Flag.Args())
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[0])
	usage()
}
