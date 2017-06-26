package commands

import (
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
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
	output, err := version.FullVersion()
	if output != "" {
		ui.Println(output)
	}
	utils.Check(err)
	args.NoForward()
}
