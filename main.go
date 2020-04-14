// +build go1.8

package main

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/github/hub/commands"
	"github.com/github/hub/github"
	"github.com/github/hub/ui"
)

func main() {
	defer github.CaptureCrash()
	err := commands.CmdRunner.Execute(os.Args)
	exitCode := handleError(err)
	os.Exit(exitCode)
}

func handleError(err error) int {
	if err == nil {
		return 0
	}

	switch e := err.(type) {
	case *exec.ExitError:
		if status, ok := e.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
		return 1
	case *commands.ErrHelp:
		ui.Println(err)
		return 0
	default:
		if errString := err.Error(); errString != "" {
			ui.Errorln(err)
		}
		return 1
	}
}
