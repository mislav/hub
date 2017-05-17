// +build go1.8

package main

import (
	"os"

	"github.com/github/hub/commands"
	"github.com/github/hub/github"
	"github.com/github/hub/ui"
)

func main() {
	defer github.CaptureCrash()

	err := commands.CmdRunner.Execute()
	if !err.Ran {
		ui.Errorln(err.Error())
	}
	os.Exit(err.ExitCode)
}
