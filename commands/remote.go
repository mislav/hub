package commands

import (
	"fmt"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"os"
	"path/filepath"
	"regexp"
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
	if args.Size() >= 2 && (args.First() == "add" || args.First() == "set-url") {
		transformRemoteArgs(args)
	}
}

func transformRemoteArgs(args *Args) {
	ownerWithName := args.Last()
	owner, name, match := parseRepoNameOwner(ownerWithName)
	if !match {
		return
	}
	isPriavte := parseRemotePrivateFlag(args)
	if name == "" {
		dir, err := os.Getwd()
		utils.Check(err)
		name = filepath.Base(dir)
	}

	if owner == "origin" {
		owner = github.CurrentConfig().FetchUser()
	}

	project := github.Project{Owner: owner, Name: name}
	url := project.GitURL(name, owner, isPriavte)

	args.Append(url)
}

func parseRepoNameOwner(nameWithOwner string) (string, string, bool) {
	ownerRe := fmt.Sprintf("^(%s)$", OwnerRe)
	ownerRegexp := regexp.MustCompile(ownerRe)
	if ownerRegexp.MatchString(nameWithOwner) {
		return ownerRegexp.FindStringSubmatch(nameWithOwner)[1], "", true
	}

	nameWithOwnerRe := fmt.Sprintf("^(%s)\\/(%s)$", OwnerRe, NameRe)
	nameWithOwnerRegexp := regexp.MustCompile(nameWithOwnerRe)
	if nameWithOwnerRegexp.MatchString(nameWithOwner) {
		match := nameWithOwnerRegexp.FindStringSubmatch(nameWithOwner)
		return match[1], match[2], true
	}

	return "", "", false
}

func parseRemotePrivateFlag(args *Args) bool {
	if i := args.IndexOf("-p"); i != -1 {
		args.Remove(i)
		return true
	}

	return false
}
