package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/github/hub/cmd"
	"github.com/github/hub/git"
	"github.com/github/hub/ui"
	"github.com/kballard/go-shellquote"
	flag "github.com/ogier/pflag"
)

type ExecError struct {
	Err      error
	Ran      bool
	ExitCode int
}

func (execError *ExecError) Error() string {
	return execError.Err.Error()
}

func newExecError(err error) ExecError {
	exitCode := 0
	ran := true

	if err != nil {
		exitCode = 1
		switch e := err.(type) {
		case *exec.ExitError:
			if status, ok := e.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		case *exec.Error:
			ran = false
		}
	}

	return ExecError{
		Err:      err,
		Ran:      ran,
		ExitCode: exitCode,
	}
}

type Runner struct {
	commands map[string]*Command
	execute  func([]*cmd.Cmd, bool) error
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
	args.ProgramPath = os.Args[0]
	forceFail := false

	if args.Command == "" && len(args.GlobalFlags) == 0 {
		args.Command = "help"
		forceFail = true
	}

	git.GlobalFlags = args.GlobalFlags // preserve git global flags
	if !isBuiltInHubCommand(args.Command) {
		expandAlias(args)
	}

	cmd := r.Lookup(args.Command)
	if cmd != nil && cmd.Runnable() {
		execErr := r.Call(cmd, args)
		if execErr.ExitCode == 0 && forceFail {
			execErr = newExecError(fmt.Errorf(""))
		}
		return execErr
	}

	gitArgs := []string{args.Command}
	gitArgs = append(gitArgs, args.Params...)

	err := git.Run(gitArgs...)
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
		err = r.execute(cmds, len(args.Callbacks) == 0)
	}

	if err == nil {
		for _, fn := range args.Callbacks {
			err = fn()
			if err != nil {
				break
			}
		}
	}

	return newExecError(err)
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
	for hubCommand, _ := range CmdRunner.All() {
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
