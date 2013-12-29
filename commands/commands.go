package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
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

	subCommands map[string]*Command
}

func (c *Command) Call(args *Args) (err error) {
	runCommand := c
	if len(c.subCommands) > 0 && args.HasSubcommand() {
		subCommandName := args.FirstParam()
		if subCommand, ok := c.subCommands[subCommandName]; ok {
			runCommand = subCommand
			args.Params = args.Params[1:]
		} else {
			fmt.Printf("error: Unknown subcommand: %s\n", subCommandName)
			c.printShortCommands()
			os.Exit(1)
		}
	}

	if err = c.parseArguments(args); err != nil {
		return
	}

	runCommand.Run(runCommand, args)
	return
}

func (c *Command) parseArguments(args *Args) (err error) {
	if err := c.Flag.Parse(args.Params); err != nil {
		if err == flag.ErrHelp {
			return nil
		} else {
			return err
		}
	}

	args.Params = c.Flag.Args()
	return
}

func (c *Command) Use(name string, subCommand *Command) {
	if c.subCommands == nil {
		c.subCommands = make(map[string]*Command)
	}
	c.subCommands[name] = subCommand
}

func (c *Command) PrintUsage() {
	if c.Runnable() {
		fmt.Printf("usage: %s\n\n", c.FormattedUsage())
	}

	fmt.Println(strings.Trim(c.Long, "\n"))
}

func (c *Command) FormattedUsage() string {
	return fmt.Sprintf("%s %s", execName(), c.Usage)
}

func (c *Command) printShortCommands() {
	if c.Runnable() {
		fmt.Printf("usage: %s\n", c.FormattedUsage())
	}
	if c.subCommands != nil && len(c.subCommands) > 0 {
		for _, s := range c.subCommands {
			fmt.Printf("   or: %s\n", s.FormattedUsage())
		}
	}
	fmt.Println()
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
	cmdApply,
	cmdCherryPick,
}

var Remote = []*Command{
	cmdClone,
	cmdFetch,
	cmdPush,
	cmdRemote,
	cmdSubmodule,
}

var GitHub = []*Command{
	cmdPullRequest,
	cmdFork,
	cmdCreate,
	cmdCiStatus,
	cmdBrowse,
	cmdCompare,
	cmdRelease,
	cmdIssue,
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
	all = append(all, cmdUpdate)

	return all
}

func execName() string {
	return filepath.Base(os.Args[0])
}
