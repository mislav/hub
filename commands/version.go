package commands

import (
	"github.com/github/hub/ui"
	"github.com/github/hub/version"
)

var cmdVersion = &Command{
	Run:          runVersion,
	Usage:        "version",
	Long:         "Shows git version and hub client version.",
	GitExtension: true,
}

func init() {
	CmdRunner.Use(cmdVersion, "--version")
}

func runVersion(cmd *Command, args *Args) {
	ui.Println(version.FullVersion())
	args.NoForward()
}
