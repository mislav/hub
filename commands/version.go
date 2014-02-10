package commands

import (
	"fmt"
	"github.com/github/hub/git"
	"github.com/github/hub/utils"
	"os"
)

const Version = "1.0.0"

var cmdVersion = &Command{
	Run:   runVersion,
	Usage: "version",
	Short: "Show gh version",
	Long:  `Shows git version and gh client version.`,
}

func init() {
	CmdRunner.Use(cmdVersion)
}

func runVersion(cmd *Command, args *Args) {
	gitVersion, err := git.Version()
	utils.Check(err)

	ghVersion := fmt.Sprintf("gh version %s", Version)

	fmt.Println(gitVersion)
	fmt.Println(ghVersion)

	os.Exit(0)
}
