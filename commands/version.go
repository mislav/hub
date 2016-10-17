package commands

import (
	"github.com/boris-rea/hub/version"
	"github.com/github/hub/ui"
)

var cmdVersion = &Command{
	Run:   runVersion,
	Usage: "version",
	Long:  "Shows git version and hub client version.",
}

func init() {
	CmdRunner.Use(cmdVersion, "--version")
}

func runVersion(cmd *Command, args *Args) {
	ui.Println(version.FullVersion())
	args.NoForward()
}
