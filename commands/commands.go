package commands

import (
	"flag"
	"fmt"
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
	if c.Runnable() {
		fmt.Printf("Usage: git %s\n\n", c.Usage)
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

var Basic = []*Command{
	cmdInit,
}

var Branching = []*Command{
	cmdCheckout,
	cmdMerge,
}

var Remote = []*Command{
	cmdClone,
	cmdFetch,
	cmdRemote,
}

var GitHub = []*Command{
	cmdPullRequest,
	cmdFork,
	cmdCreate,
	cmdCiStatus,
	cmdBrowse,
	cmdCompare,
}

func All() []*Command {
	all := make([]*Command, 0)
	all = append(all, Basic...)
	all = append(all, Branching...)
	all = append(all, Remote...)
	all = append(all, GitHub...)
	all = append(all, cmdAlias)
	all = append(all, cmdVersion)
	all = append(all, cmdHelp)

	return all
}
