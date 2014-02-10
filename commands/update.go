package commands

import (
	"github.com/github/hub/utils"
	"os"
)

var cmdUpdate = &Command{
	Run:   update,
	Usage: "update",
	Short: "Update gh",
	Long: `Update gh to the latest version.

Examples:
  git update
`,
}

func init() {
	CmdRunner.Use(cmdUpdate)
}

func update(cmd *Command, args *Args) {
	updater := NewUpdater()
	err := updater.Update()
	utils.Check(err)
	os.Exit(0)
}
