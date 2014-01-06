package main

import (
	"github.com/jingweno/gh/commands"
	"github.com/jingweno/gh/github"
	"os"
)

func main() {
	defer github.CaptureCrash()

	err := commands.CmdRunner.Execute()
	os.Exit(err.ExitCode)
}
