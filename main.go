package main

import (
	"github.com/jingweno/gh/commands"
	"os"
)

func main() {
	runner := commands.Runner{Args: os.Args[1:]}
	err := runner.Execute()
	os.Exit(err.ExitCode)
}
