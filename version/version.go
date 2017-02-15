package version

import (
	"fmt"

	"github.com/github/hub/git"
	"github.com/github/hub/utils"
)

var Version = "2.3.0-pre9"

func FullVersion() string {
	gitVersion, err := git.Version()
	utils.Check(err)
	return fmt.Sprintf("%s\nhub version %s", gitVersion, Version)
}
