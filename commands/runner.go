package commands

import (
	"fmt"
	"strings"

	"github.com/github/hub/v2/cmd"
	"github.com/github/hub/v2/git"
	"github.com/github/hub/v2/ui"
	"github.com/kballard/go-shellquote"
)

type Runner struct {
	commands map[string]*Command
}

func NewRunner() *Runner {
	return &Runner{
		commands: make(map[string]*Command),
	}
}

func (r *Runner) All() map[string]*Command {
	return r.commands
}

func (r *Runner) Use(command *Command, aliases ...string) {
	r.commands[command.Name()] = command
	if len(aliases) > 0 {
		r.commands[aliases[0]] = command
	}
}

func (r *Runner) Lookup(name string) *Command {
	return r.commands[name]
}

func (r *Runner) Execute(cliArgs []string) error {
	args := NewArgs(cliArgs[1:])
	args.ProgramPath = cliArgs[0]
	forceFail := false

	if args.Command == "" && len(args.GlobalFlags) == 0 {
		args.Command = "help"
		forceFail = true
	}

	cmdName := args.Command
	if strings.Contains(cmdName, "=") {
		cmdName = strings.SplitN(cmdName, "=", 2)[0]
	}

	git.GlobalFlags = args.GlobalFlags // preserve git global flags
	if !isBuiltInHubCommand(cmdName) {
		expandAlias(args)
		cmdName = args.Command
	}

	// make `<cmd> --help` equivalent to `help <cmd>`
	if args.ParamsSize() == 1 && args.GetParam(0) == helpFlag {
		if c := r.Lookup(cmdName); c != nil && !c.GitExtension {
			args.ReplaceParam(0, cmdName)
			args.Command = "help"
			cmdName = args.Command
		}
	}

	cmd := r.Lookup(cmdName)
	if cmd != nil && cmd.Runnable() {
		err := callRunnableCommand(cmd, args)
		if err == nil && forceFail {
			err = fmt.Errorf("")
		}
		return err
	}

	gitArgs := []string{}
	if args.Command != "" {
		gitArgs = append(gitArgs, args.Command)
	}
	gitArgs = append(gitArgs, args.Params...)

	return git.Run(gitArgs...)
}

func callRunnableCommand(cmd *Command, args *Args) error {
	err := cmd.Call(args)
	if err != nil {
		return err
	}

	cmds := args.Commands()
	if args.Noop {
		printCommands(cmds)
	} else if err = executeCommands(cmds, len(args.Callbacks) == 0); err != nil {
		return err
	}

	for _, fn := range args.Callbacks {
		if err = fn(); err != nil {
			return err
		}
	}

	return nil
}

func printCommands(cmds []*cmd.Cmd) {
	for _, c := range cmds {
		ui.Println(c)
	}
}

func executeCommands(cmds []*cmd.Cmd, execFinal bool) error {
	for i, c := range cmds {
		var err error
		// Run with `Exec` for the last command in chain
		if execFinal && i == len(cmds)-1 {
			err = c.Run()
		} else {
			err = c.Spawn()
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func expandAlias(args *Args) {
	cmd := args.Command
	if cmd == "" {
		return
	}
	expandedCmd, err := git.Alias(cmd)

	if err == nil && expandedCmd != "" && !git.IsBuiltInGitCommand(cmd) {
		words, e := splitAliasCmd(expandedCmd)
		if e == nil {
			args.Command = words[0]
			args.PrependParams(words[1:]...)
		}
	}
}

func isBuiltInHubCommand(command string) bool {
	for hubCommand := range CmdRunner.All() {
		if hubCommand == command {
			return true
		}
	}
	return false
}

func splitAliasCmd(cmd string) ([]string, error) {
	if cmd == "" {
		return nil, fmt.Errorf("alias can't be empty")
	}

	if strings.HasPrefix(cmd, "!") {
		return nil, fmt.Errorf("alias starting with ! can't be split")
	}

	words, err := shellquote.Split(cmd)
	if err != nil {
		return nil, err
	}

	return words, nil
}
