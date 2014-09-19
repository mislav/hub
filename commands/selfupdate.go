package commands

import (
	"os"

	"github.com/github/hub/utils"
)

var cmdSelfupdate = &Command{
	Run:   update,
	Usage: "selfupdate",
	Short: "Update Hub",
	Long: `Update Hub to the latest version.

Examples:
  git selfupdate
`,
}

func init() {
	CmdRunner.Use(cmdSelfupdate)
}

func update(cmd *Command, args *Args) {
	updater := NewUpdater()
	err := updater.Update()
	utils.Check(err)
	os.Exit(0)
}
