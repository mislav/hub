package commands

import (
	"fmt"
	"github.com/jingweno/gh/cmd"
	"github.com/jingweno/gh/git"
)

type Runner struct {
	Args []string
}

func (r *Runner) Execute() error {
	args := NewArgs(r.Args)
	if args.Command == "" {
		usage()
	}

	expandAlias(args)
	slurpGlobalFlags(args)

	for _, cmd := range All() {
		if cmd.Name() == args.Command && cmd.Runnable() {
			if !cmd.GitExtension {
				cmd.Flag.Usage = func() {
					cmd.PrintUsage()
				}
				if err := cmd.Flag.Parse(args.Params); err != nil {
					return err
				}

				args.Params = cmd.Flag.Args()
			}

			cmd.Run(cmd, args)

			cmds := args.Commands()
			if args.Noop {
				printCommands(cmds)
			} else {
				err := executeCommands(cmds)
				if err != nil {
					return err
				}
			}

			return nil
		}
	}

	return git.SysExec(args.Command, args.Params...)
}

func slurpGlobalFlags(args *Args) {
	for i, p := range args.Params {
		if p == "--no-op" {
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
	length := len(cmds)
	for i, c := range cmds {
		var err error
		if i == (length - 1) {
			err = c.SysExec()
		} else {
			err = c.Exec()
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func expandAlias(args *Args) {
	cmd := args.Command
	expandedCmd, err := git.Config(fmt.Sprintf("alias.%s", cmd))
	if err == nil && expandedCmd != "" {
		args.Command = expandedCmd
	}
}
