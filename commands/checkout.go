package commands

import (
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/utils"
)

var cmdCheckout = &Command{
	Run:          checkout,
	GitExtension: true,
	Usage:        "checkout PULLREQ-URL [BRANCH]",
	Short:        "Switch the active branch to another branch",
}

/**
  $ gh checkout https://github.com/jingweno/gh/pull/73
  # > git remote add -f -t feature git://github:com/foo/gh.git
  # > git checkout --track -B foo-feature foo/feature

  $ gh checkout https://github.com/jingweno/gh/pull/73 custom-branch-name
**/
func checkout(command *Command, args []string) {
	if len(args) > 0 {
		args = transformCheckoutArgs(args)
	}

  err := git.ExecCheckout(args)
	utils.Check(err)
}

func transformCheckoutArgs(args []string) []string {
  return nil
}
