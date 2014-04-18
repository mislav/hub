package commands

import (
	"github.com/github/hub/utils"
	"os"
)

var cmdSelfupdate = &Command{
	Run:   update,
	Usage: "selfupdate",
	Short: "Update gh",
	Long: `Update gh to the latest version.

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
