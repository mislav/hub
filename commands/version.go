package commands

import (
	"os"

	"github.com/github/hub/ui"
	"github.com/github/hub/version"
)

var cmdVersion = &Command{
	Run:   runVersion,
	Usage: "version",
	Short: "Show hub version",
	Long:  `Shows git version and hub client version.`,
}

func init() {
	CmdRunner.Use(cmdVersion, "--version")
}

func runVersion(cmd *Command, args *Args) {
	ui.Println(version.FullVersion())
	os.Exit(0)
}
