package commands

import (
	"flag"
	"fmt"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/utils"
	"strings"
)

var (
	NameRe          = "[\\w.][\\w.-]*"
	OwnerRe         = "[a-zA-Z0-9][a-zA-Z0-9-]*"
	NameWithOwnerRe = fmt.Sprintf("^(?:%s|%s\\/%s)$", NameRe, OwnerRe, NameRe)
)

type Command struct {
	Run  func(cmd *Command, args *Args)
	Flag flag.FlagSet

	Usage        string
	Short        string
	Long         string
	GitExtension bool
}

func (c *Command) PrintUsage() {
	if c.GitExtension {
		err := git.SysExec(c.Name(), "--help")
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

var Branching = []*Command{
	cmdCheckout,
}

var Remote = []*Command{
	cmdClone,
	cmdRemote,
}

var GitHub = []*Command{
	cmdPullRequest,
	cmdFork,
	cmdCiStatus,
	cmdBrowse,
	cmdCompare,
}

func All() []*Command {
	all := make([]*Command, 0)
	all = append(all, Branching...)
	all = append(all, Remote...)
	all = append(all, GitHub...)
	all = append(all, cmdAlias)
	all = append(all, cmdVersion)
	all = append(all, cmdHelp)

	return all
}
