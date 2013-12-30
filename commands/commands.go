package commands

import (
	"bytes"
	"fmt"
	flag "github.com/ogier/pflag"
	"os"
	"path/filepath"
	"strings"
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
	runCommand, err := lookupCommand(c, args)
	if err != nil {
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
	if !c.GitExtension {
		c.Flag.Usage = c.PrintUsage
	}

	if err = c.Flag.Parse(args.Params); err == nil {
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
		fmt.Printf("usage: %s\n\n", c.FormattedUsage())
	}

	fmt.Println(strings.Trim(c.Long, "\n"))
}

func (c *Command) FormattedUsage() string {
	return fmt.Sprintf("%s %s", execName(), c.Usage)
}

func (c *Command) subCommandsUsage() string {
	buffer := bytes.NewBufferString("")

	key := "usage"

	key = printUsageBuffer(c, buffer, key)
	if c.subCommands != nil && len(c.subCommands) > 0 {
		for _, s := range c.subCommands {
			key = printUsageBuffer(s, buffer, key)
		}
	}
	buffer.WriteString("\n")

	return buffer.String()
}

func printUsageBuffer(c *Command, b *bytes.Buffer, key string) string {
	if c.Runnable() {
		b.WriteString(fmt.Sprintf("%s: %s\n", key, c.FormattedUsage()))
		key = "   or"
	}
	return key
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

func execName() string {
	return filepath.Base(os.Args[0])
}

func lookupCommand(c *Command, args *Args) (runCommand *Command, err error) {
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
