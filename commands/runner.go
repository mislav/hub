package commands

import (
	"fmt"
	"github.com/jingweno/gh/cmd"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/utils"
	"github.com/kballard/go-shellquote"
	flag "github.com/ogier/pflag"
	"os"
	"os/exec"
	"strings"
	"syscall"
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
}

func NewRunner() *Runner {
	return &Runner{commands: make(map[string]*Command)}
}

func (r *Runner) All() map[string]*Command {
	return r.commands
}

func (r *Runner) Use(command *Command) {
	r.commands[command.Name()] = command
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

	expandAlias(args)
	slurpGlobalFlags(args)

	cmd := r.Lookup(args.Command)
	if cmd != nil && cmd.Runnable() {
		return r.Call(cmd, args)
	}

	err = git.Spawn(args.Command, args.Params...)
	return newExecError(err)
}

func (r *Runner) Call(cmd *Command, args *Args) ExecError {
	if !cmd.GitExtension {
		cmd.Flag.Usage = func() {
			cmd.PrintUsage()
		}
	}

	err := cmd.Call(args)
	if err != nil && err != flag.ErrHelp {
		return newExecError(err)
	}

	cmds := args.Commands()
	if args.Noop {
		printCommands(cmds)
	} else {
		err = executeCommands(cmds)
	}

	return newExecError(err)
}

func slurpGlobalFlags(args *Args) {
	for i, p := range args.Params {
		if p == "--noop" {
			args.Noop = true
			args.RemoveParam(i)
		}
	}
}

func printCommands(cmds []*cmd.Cmd) {
	for _, c := range cmds {
		fmt.Println(c)
	}
}

func executeCommands(cmds []*cmd.Cmd) error {
	for _, c := range cmds {
		err := c.Exec()
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
