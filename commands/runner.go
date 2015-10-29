package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/github/hub/Godeps/_workspace/src/github.com/kballard/go-shellquote"
	flag "github.com/github/hub/Godeps/_workspace/src/github.com/ogier/pflag"
	"github.com/github/hub/cmd"
	"github.com/github/hub/git"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

type ExecError struct {
	Err      error
	ExitCode int
}

func (execError *ExecError) Error() string {
	return execError.Err.Error()
}

func newExecError(err error) ExecError {
	exitCode := 0
	if err != nil {
		exitCode = 1
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		}
	}

	return ExecError{Err: err, ExitCode: exitCode}
}

type Runner struct {
	commands map[string]*Command
	execute  func(cmds []*cmd.Cmd) error
}

func NewRunner() *Runner {
	return &Runner{
		commands: make(map[string]*Command),
		execute:  executeCommands,
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

func (r *Runner) Execute() ExecError {
	args := NewArgs(os.Args[1:])

	if args.Command == "" {
		printUsage()
		return newExecError(nil)
	}

	updater := NewUpdater()
	err := updater.PromptForUpdate()
	utils.Check(err)

	git.GlobalFlags = args.GlobalFlags // preserve git global flags
	expandAlias(args)

	cmd := r.Lookup(args.Command)
	if cmd != nil && cmd.Runnable() {
		return r.Call(cmd, args)
	}

	err = git.Run(args.Command, args.Params...)
	return newExecError(err)
}

func (r *Runner) Call(cmd *Command, args *Args) ExecError {
	err := cmd.Call(args)
	if err != nil {
		if err == flag.ErrHelp {
			err = nil
		}
		return newExecError(err)
	}

	cmds := args.Commands()
	if args.Noop {
		printCommands(cmds)
	} else {
		err = r.execute(cmds)
	}

	return newExecError(err)
}

func printCommands(cmds []*cmd.Cmd) {
	for _, c := range cmds {
		ui.Println(c)
	}
}

func executeCommands(cmds []*cmd.Cmd) error {
	for i, c := range cmds {
		var err error
		// Run with `Exec` for the last command in chain
		if i == len(cmds)-1 {
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
	expandedCmd, err := git.Alias(cmd)
	if err == nil && expandedCmd != "" {
		words, e := splitAliasCmd(expandedCmd)
		if e == nil {
			args.Command = words[0]
			args.PrependParams(words[1:]...)
		}
	}
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
