package commands

import (
	"bytes"
	"fmt"
	"strings"

	flag "github.com/github/hub/Godeps/_workspace/src/github.com/ogier/pflag"
	"github.com/github/hub/ui"
)

var (
	NameRe          = "[\\w.][\\w.-]*"
	OwnerRe         = "[a-zA-Z0-9][a-zA-Z0-9-]*"
	NameWithOwnerRe = fmt.Sprintf("^(?:%s|%s\\/%s)$", NameRe, OwnerRe, NameRe)

	CmdRunner = NewRunner()
)

type Command struct {
	Run  func(cmd *Command, args *Args)
	Flag flag.FlagSet

	Key          string
	Usage        string
	Short        string
	Long         string
	GitExtension bool

	subCommands map[string]*Command
}

func (c *Command) Call(args *Args) (err error) {
	runCommand, err := c.lookupSubCommand(args)
	if err != nil {
		ui.Errorln(err)
		return
	}

	if !c.GitExtension {
		err = runCommand.parseArguments(args)
		if err != nil {
			return
		}
	}

	runCommand.Run(runCommand, args)

	return
}

func (c *Command) parseArguments(args *Args) (err error) {
	c.Flag.SetInterspersed(true)
	c.Flag.Init(c.Name(), flag.ContinueOnError)
	c.Flag.Usage = c.PrintUsage
	if err = c.Flag.Parse(args.Params); err == nil {
		for _, arg := range args.Params {
			if arg == "--" {
				args.Terminator = true
			}
		}
		args.Params = c.Flag.Args()
	}

	return
}

func (c *Command) Use(subCommand *Command) {
	if c.subCommands == nil {
		c.subCommands = make(map[string]*Command)
	}
	c.subCommands[subCommand.Name()] = subCommand
}

func (c *Command) PrintUsage() {
	if c.Runnable() {
		ui.Printf("usage: %s\n\n", c.FormattedUsage())
	}

	ui.Println(strings.Trim(c.Long, "\n"))
}

func (c *Command) FormattedUsage() string {
	return fmt.Sprintf("git %s", c.Usage)
}

func (c *Command) subCommandsUsage() string {
	buffer := bytes.NewBufferString("")

	usage := "usage"
	usage = printUsageBuffer(c, buffer, usage)
	for _, s := range c.subCommands {
		usage = printUsageBuffer(s, buffer, usage)
	}

	return buffer.String()
}

func printUsageBuffer(c *Command, b *bytes.Buffer, usage string) string {
	if c.Runnable() {
		b.WriteString(fmt.Sprintf("%s: %s\n", usage, c.FormattedUsage()))
		usage = "   or"
	}
	return usage
}

func (c *Command) Name() string {
	if c.Key != "" {
		return c.Key
	}
	return strings.Split(c.Usage, " ")[0]
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}

func (c *Command) List() bool {
	return c.Short != ""
}

func (c *Command) lookupSubCommand(args *Args) (runCommand *Command, err error) {
	if len(c.subCommands) > 0 && args.HasSubcommand() {
		subCommandName := args.FirstParam()
		if subCommand, ok := c.subCommands[subCommandName]; ok {
			runCommand = subCommand
			args.Params = args.Params[1:]
		} else {
			err = fmt.Errorf("error: Unknown subcommand: %s\n%s", subCommandName, c.subCommandsUsage())
		}
	} else {
		runCommand = c
	}

	return
}
