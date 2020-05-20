package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/github/hub/v2/utils"
)

var (
	NameRe          = `[\w.-]+`
	OwnerRe         = "[a-zA-Z0-9][a-zA-Z0-9-]*"
	NameWithOwnerRe = fmt.Sprintf(`^(%s/)?%s$`, OwnerRe, NameRe)

	CmdRunner = NewRunner()
)

type Command struct {
	Run func(cmd *Command, args *Args)

	Key          string
	Usage        string
	Long         string
	KnownFlags   string
	GitExtension bool

	subCommands   map[string]*Command
	parentCommand *Command
}

func (c *Command) Call(args *Args) (err error) {
	runCommand, err := c.lookupSubCommand(args)
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

type ErrHelp struct {
	err string
}

func (e ErrHelp) Error() string {
	return e.err
}

func (c *Command) parseArguments(args *Args) error {
	knownFlags := c.KnownFlags
	if knownFlags == "" {
		knownFlags = c.Long
	}
	args.Flag = utils.NewArgsParserWithUsage("-h, --help\n" + knownFlags)

	rest, err := args.Flag.Parse(args.Params)
	if err != nil {
		return fmt.Errorf("%s\n%s", err, c.Synopsis())
	}
	if args.Flag.Bool("--help") {
		return &ErrHelp{err: c.Synopsis()}
	}
	args.Params = rest
	args.Terminator = args.Flag.HasTerminated
	return nil
}

func (c *Command) Use(subCommand *Command) {
	if c.subCommands == nil {
		c.subCommands = make(map[string]*Command)
	}
	c.subCommands[subCommand.Name()] = subCommand
	subCommand.parentCommand = c
}

func (c *Command) UsageError(msg string) error {
	nl := ""
	if msg != "" {
		nl = "\n"
	}
	return fmt.Errorf("%s%s%s", msg, nl, c.Synopsis())
}

func (c *Command) Synopsis() string {
	lines := []string{}
	usagePrefix := "Usage:"
	usageStr := c.Usage
	if usageStr == "" && c.parentCommand != nil {
		usageStr = c.parentCommand.Usage
	}

	for _, line := range strings.Split(usageStr, "\n") {
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

	long = strings.Replace(long, "''", "`", -1)
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
