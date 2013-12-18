package commands

import (
	"os"
)

var cmdUpdate = &Command{
	Run:   update,
	Usage: "update",
	Short: "Update gh",
	Long: `Update gh with the latest version.

Examples:
  git update
`,
}

func update(cmd *Command, args *Args) {
	os.Exit(0)
}
