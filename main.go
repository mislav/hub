package main

import (
	"github.com/github/hub/commands"
	"github.com/github/hub/github"
	"os"
)

func main() {
	defer github.CaptureCrash()

	err := commands.CmdRunner.Execute()
	os.Exit(err.ExitCode)
}
