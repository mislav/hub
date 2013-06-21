package commands

import (
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
)

var cmdRemote = &Command{
	Run:          remote,
	GitExtension: true,
	Usage:        "remote [-p] OPTIONS USER[/REPOSITORY]",
	Short:        "View and manage a set of remote repositories",
}

/**
  $ gh remote add jingweno
  > git remote add jingweno git://github.com/jingweno/THIS_REPO.git

  $ gh remote add -p jingweno
  > git remote add jingweno git@github.com:jingweno/THIS_REPO.git

  $ gh remote add origin
  > git remote add origin
  git://github.com/YOUR_LOGIN/THIS_REPO.git
**/
func remote(command *Command, args []string) {
	if len(args) >= 2 && (args[0] == "add" || args[0] == "set-url") {
		args = transformRemoteArgs(args)
	}

	err := git.ExecRemote(args...)
	utils.Check(err)
}

func transformRemoteArgs(args []string) (newArgs []string) {
	args, isPriavte := parseRemotePrivateFlag(args)
	newArgs, owner := removeItem(args, len(args)-1)

	gh := github.New()
	url := gh.ExpandRemoteUrl(owner, isPriavte)

	return append(newArgs, owner, url)
}

func parseRemotePrivateFlag(args []string) ([]string, bool) {
	for i, arg := range args {
		if arg == "-p" {
			args, _ = removeItem(args, i)
			return args, true
		}
	}

	return args, false
}
