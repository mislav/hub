package commands

import (
	"os"

	"github.com/github/hub/utils"
)

var cmdSelfupdate = &Command{
	Run:   update,
	Usage: "selfupdate",
	Long:  "Update hub to the latest version.",
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
