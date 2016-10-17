package main

import (
	"os"

	"github.com/boris-rea/hub/commands"
	"github.com/boris-rea/hub/github"
	"github.com/boris-rea/hub/ui"
)

func main() {
	defer github.CaptureCrash()

	err := commands.CmdRunner.Execute()
	if !err.Ran {
		ui.Errorln(err.Error())
	}
	os.Exit(err.ExitCode)
}
