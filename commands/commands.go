package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/github/hub/ui"
	flag "github.com/ogier/pflag"
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
	c.Flag.Usage = func() {
		if args.HasFlags("-help", "--help") {
			ui.Println(c.Synopsis())
		} else {
			ui.Errorln(c.Synopsis())
		}
	}
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

func (c *Command) FlagPassed(name string) bool {
	found := false
	c.Flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func (c *Command) Arg(idx int) string {
	args := c.Flag.Args()
	if idx < len(args) {
		return args[idx]
	} else {
		return ""
	}
}

func (c *Command) Use(subCommand *Command) {
	if c.subCommands == nil {
		c.subCommands = make(map[string]*Command)
	}
	c.subCommands[subCommand.Name()] = subCommand
}

func (c *Command) Synopsis() string {
	lines := []string{}
	usagePrefix := "Usage:"

	for _, line := range strings.Split(c.Usage, "\n") {
		if line != "" {
			usage := fmt.Sprintf("%s hub %s", usagePrefix, line)
			usagePrefix = "      "
			lines = append(lines, usage)
		}
	}
	return strings.Join(lines, "\n")
}

func (c *Command) HelpText() string {
	usage := strings.Replace(c.Usage, "-^", "`-^`", 1)
	usageRe := regexp.MustCompile(`(?m)^([a-z-]+)(.*)$`)
	usage = usageRe.ReplaceAllString(usage, "`hub $1`$2  ")
	usage = strings.TrimSpace(usage)

	var desc string
	long := strings.TrimSpace(c.Long)
	if lines := strings.Split(long, "\n"); len(lines) > 1 {
		desc = lines[0]
		long = strings.Join(lines[1:], "\n")
	}

	long = strings.Replace(long, "'", "`", -1)
	headingRe := regexp.MustCompile(`(?m)^(## .+):$`)
	long = headingRe.ReplaceAllString(long, "$1")

	indentRe := regexp.MustCompile(`(?m)^\t`)
	long = indentRe.ReplaceAllLiteralString(long, "")
	definitionListRe := regexp.MustCompile(`(?m)^(\* )?([^#\s][^\n]*?):?\n\t`)
	long = definitionListRe.ReplaceAllString(long, "$2\n:\t")

	return fmt.Sprintf("hub-%s(1) -- %s\n===\n\n## Synopsis\n\n%s\n%s", c.Name(), desc, usage, long)
}

func (c *Command) Name() string {
	if c.Key != "" {
		return c.Key
	}
	usageLine := strings.Split(strings.TrimSpace(c.Usage), "\n")[0]
	return strings.Split(usageLine, " ")[0]
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}

func (c *Command) lookupSubCommand(args *Args) (runCommand *Command, err error) {
	if len(c.subCommands) > 0 && args.HasSubcommand() {
		subCommandName := args.FirstParam()
		if subCommand, ok := c.subCommands[subCommandName]; ok {
			runCommand = subCommand
			args.Params = args.Params[1:]
		} else {
			err = fmt.Errorf("error: Unknown subcommand: %s", subCommandName)
		}
	} else {
		runCommand = c
	}

	return
}
