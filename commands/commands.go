package commands

import (
	"flag"
	"fmt"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/utils"
	"strings"
)

type Command struct {
	Run  func(cmd *Command, args []string)
	Flag flag.FlagSet

	Usage        string
	Short        string
	Long         string
	GitExtension bool
}

func (c *Command) PrintUsage() {
	if c.GitExtension {
		err := git.Help(c.Name())
		utils.Check(err)
	} else {
		if c.Runnable() {
			fmt.Printf("Usage: gh %s\n\n", c.Usage)
		}

		fmt.Println(strings.Trim(c.Long, "\n"))
	}
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

var All = append(Remote, GitHub...)

var Remote = []*Command{
	cmdRemote,
}

var GitHub = []*Command{
	cmdPull,
	cmdFork,
	cmdCi,
	cmdBrowse,
	cmdCompare,
	cmdHelp,
	cmdVersion,
}
