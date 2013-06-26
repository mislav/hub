package main

import (
	"fmt"
	"github.com/jingweno/gh/commands"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/utils"
	"os"
)

func main() {
	args := commands.NewArgs(os.Args[1:])
	if args.Size() < 1 {
		commands.Usage()
		return
	}

	for _, cmd := range commands.All() {
		if cmd.Name() == args.First() && cmd.Runnable() {
			cmdArgs := args.Rest()
			if !cmd.GitExtension {
				cmd.Flag.Usage = func() {
					cmd.PrintUsage()
				}
				if err := cmd.Flag.Parse(cmdArgs); err != nil {
					os.Exit(2)
				}

				cmdArgs = cmd.Flag.Args()
			}

			args = commands.NewArgs(cmdArgs)
			cmd.Run(cmd, args)
			return
		}
	}

	if args.Size() > 0 {
		err := git.SysExec(args.First(), args.Rest()...)
		utils.Check(err)
		return
	}

	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args.First())
	commands.Usage()
}
