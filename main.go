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

	for _, cmd := range commands.All {
		if cmd.Name() == args[0] && cmd.Runnable() {
			cmd.Flag.Usage = func() {
				cmd.PrintUsage()
			}
			if err := cmd.Flag.Parse(args[1:]); err != nil {
				os.Exit(2)
			}
			cmd.Run(cmd, cmd.Flag.Args())
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[0])
	commands.Usage()
}
