package main

import (
	"fmt"
	"github.com/jingweno/gh/commands"
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

	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[0])
	commands.Usage()
}
