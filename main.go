package main

import (
	"fmt"
	"github.com/jingweno/gh/commands"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/utils"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		commands.Usage()
	}

	for _, cmd := range commands.All() {
		if cmd.Name() == args[0] && cmd.Runnable() {
			cmdArgs := args[1:]
			if !cmd.GitExtension {
				cmd.Flag.Usage = func() {
					cmd.PrintUsage()
				}
				if err := cmd.Flag.Parse(args[1:]); err != nil {
					os.Exit(2)
				}

				cmdArgs = cmd.Flag.Args()
			}

			cmd.Run(cmd, cmdArgs)

			return
		}
	}

	if len(args) > 0 {
		err := git.SysExec(args[0], args[1:]...)
		utils.Check(err)
		return
	}

	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[0])
	commands.Usage()
}
