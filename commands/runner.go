package commands

import (
	"fmt"
	"github.com/jingweno/gh/git"
	"os"
)

type Runner struct {
	Args []string
}

func (r *Runner) Execute() error {
	args := NewArgs(os.Args[1:])
	if args.Size() < 1 {
		usage()
	}

	expandAlias(args)

	for _, cmd := range All() {
		if cmd.Name() == args.First() && cmd.Runnable() {
			cmdArgs := args.Rest()
			if !cmd.GitExtension {
				cmd.Flag.Usage = func() {
					cmd.PrintUsage()
				}
				if err := cmd.Flag.Parse(cmdArgs); err != nil {
					return err
				}

				cmdArgs = cmd.Flag.Args()
			}

			args = NewArgs(cmdArgs)
			cmd.Run(cmd, args)
			args.Prepend(cmd.Name())

			cmds := args.Commands()
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
	}

	return git.SysExec(args.First(), args.Rest()...)
}

func expandAlias(args *Args) {
	cmd := args.First()
	expandedCmd, err := git.Config(fmt.Sprintf("alias.%s", cmd))
	if err == nil && expandedCmd != "" {
		args.Remove(0)
		args.Prepend(expandedCmd)
	}
}
