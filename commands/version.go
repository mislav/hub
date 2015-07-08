package commands

import (
	"fmt"
	"os"

	"github.com/github/hub/git"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var Version = "2.2.0"

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
	gitVersion, err := git.Version()
	utils.Check(err)

	ghVersion := fmt.Sprintf("hub version %s", Version)

	ui.Println(gitVersion)
	ui.Println(ghVersion)

	os.Exit(0)
}
