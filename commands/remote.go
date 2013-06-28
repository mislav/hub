package commands

import (
	"github.com/jingweno/gh/github"
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
func remote(command *Command, args *Args) {
	if args.Size() >= 2 && (args.First() == "add" || args.First() == "set-url") {
		transformRemoteArgs(args)
	}
}

func transformRemoteArgs(args *Args) {
	isPriavte := parseRemotePrivateFlag(args)
	owner := args.Last()

	gh := github.New()
	url := gh.ExpandRemoteUrl(owner, isPriavte)

	args.Append(url)
}

func parseRemotePrivateFlag(args *Args) bool {
	if i := args.IndexOf("-p"); i != -1 {
		args.Remove(i)
		return true
	}

	return false
}
