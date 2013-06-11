package commands

import (
	"fmt"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/utils"
)

const Version = "0.6.0"

var cmdVersion = &Command{
	Run:   runVersion,
	Usage: "version",
	Short: "Show gh version",
	Long:  `Shows git version and gh client version.`,
}

func runVersion(cmd *Command, args []string) {
	gitVersion, err := git.Version()
	utils.Check(err)

	ghVersion := fmt.Sprintf("gh version %s", Version)

	fmt.Println(gitVersion)
	fmt.Println(ghVersion)
}
