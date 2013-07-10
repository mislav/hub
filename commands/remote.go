package commands

import (
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
)

var cmdRemote = &Command{
	Run:          remote,
	GitExtension: true,
	Usage:        "remote [-p] OPTIONS USER[/REPOSITORY]",
	Short:        "View and manage a set of remote repositories",
	Long: `Add remote "git://github.com/USER/REPOSITORY.git" as with
git-remote(1). When /REPOSITORY is omitted, the basename of the
current working directory is used. With -p, use private remote
"git@github.com:USER/REPOSITORY.git". If USER is "origin"
then uses your GitHub login.
`,
}

/*
  $ gh remote add jingweno
  > git remote add jingweno git://github.com/jingweno/THIS_REPO.git

  $ gh remote add -p jingweno
  > git remote add jingweno git@github.com:jingweno/THIS_REPO.git

  $ gh remote add origin
  > git remote add origin git://github.com/YOUR_LOGIN/THIS_REPO.git
*/
func remote(command *Command, args *Args) {
	if args.ParamsSize() >= 2 && (args.FirstParam() == "add" || args.FirstParam() == "set-url") {
		transformRemoteArgs(args)
	}
}

func transformRemoteArgs(args *Args) {
	ownerWithName := args.LastParam()
	owner, name, match := parseRepoNameOwner(ownerWithName)
	if !match {
		return
	}
	isPriavte := parseRemotePrivateFlag(args)
	var err error
	if name == "" {
		name, err = utils.DirName()
		utils.Check(err)
	}

	if owner == "origin" {
		owner = github.CurrentConfig().FetchUser()
	}

	project := github.Project{Owner: owner, Name: name}
	url := project.GitURL(name, owner, isPriavte)

	args.AppendParams(url)
}

func parseRemotePrivateFlag(args *Args) bool {
	if i := args.IndexOfParam("-p"); i != -1 {
		args.RemoveParam(i)
		return true
	}

	return false
}
