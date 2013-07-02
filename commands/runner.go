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
	if args.Command == "" {
		usage()
	}

	expandAlias(args)

	for _, cmd := range All() {
		if cmd.Name() == args.Command && cmd.Runnable() {
			if !cmd.GitExtension {
				cmd.Flag.Usage = func() {
					cmd.PrintUsage()
				}
				cmdArgs := args.Array()
				if err := cmd.Flag.Parse(cmdArgs); err != nil {
					return err
				}

				cmdArgs = cmd.Flag.Args()
				newArgs := []string{cmd.Name()}
				newArgs = append(newArgs, cmdArgs...)
				args = NewArgs(newArgs)
			}

			cmd.Run(cmd, args)

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

	return git.SysExec(args.Command, args.Array()...)
}

func expandAlias(args *Args) {
	cmd := args.Command
	expandedCmd, err := git.Config(fmt.Sprintf("alias.%s", cmd))
	if err == nil && expandedCmd != "" {
		args.Command = expandedCmd
	}
}
