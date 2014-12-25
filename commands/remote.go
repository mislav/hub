package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
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

func init() {
	CmdRunner.Use(cmdRemote)
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
	if !args.IsParamsEmpty() && (args.FirstParam() == "add" || args.FirstParam() == "set-url") {
		transformRemoteArgs(args)
	}
}

func transformRemoteArgs(args *Args) {
	ownerWithName := args.LastParam()
	owner, name := parseRepoNameOwner(ownerWithName)
	if owner == "" {
		return
	}

	localRepo, err := github.LocalRepo()
	utils.Check(err)

	var repoName, host string
	if name == "" {
		project, err := localRepo.MainProject()
		if err == nil {
			repoName = project.Name
			host = project.Host
		} else {
			repoName, err = utils.DirName()
			utils.Check(err)
		}

		name = repoName
	}

	hostConfig, err := github.CurrentConfig().DefaultHost()
	if err != nil {
		utils.Check(github.FormatError("adding remote", err))
	}

	words := args.Words()
	isPrivate := parseRemotePrivateFlag(args)
	if len(words) == 2 && words[1] == "origin" {
		// Origin special case triggers default user/repo
		owner = hostConfig.User
		name = repoName
	} else if len(words) == 2 {
		// gh remote add jingweno foo/bar
		if idx := args.IndexOfParam(words[1]); idx != -1 {
			args.ReplaceParam(idx, owner)
		}
	} else {
		args.RemoveParam(args.ParamsSize() - 1)
	}

	if strings.ToLower(owner) == strings.ToLower(hostConfig.User) {
		owner = hostConfig.User
		isPrivate = true
	}

	project := github.NewProject(owner, name, host)
	// for GitHub Enterprise
	isPrivate = isPrivate || project.Host != github.GitHubHost
	url := project.GitURL(name, owner, isPrivate)
	args.AppendParams(url)
}

func parseRemotePrivateFlag(args *Args) bool {
	if i := args.IndexOfParam("-p"); i != -1 {
		args.RemoveParam(i)
		return true
	}

	return false
}

func parseRepoNameOwner(nameWithOwner string) (owner, name string) {
	ownerRe := fmt.Sprintf("^(%s)$", OwnerRe)
	ownerRegexp := regexp.MustCompile(ownerRe)
	if ownerRegexp.MatchString(nameWithOwner) {
		owner = ownerRegexp.FindStringSubmatch(nameWithOwner)[1]
		return
	}

	nameWithOwnerRe := fmt.Sprintf("^(%s)\\/(%s)$", OwnerRe, NameRe)
	nameWithOwnerRegexp := regexp.MustCompile(nameWithOwnerRe)
	if nameWithOwnerRegexp.MatchString(nameWithOwner) {
		result := nameWithOwnerRegexp.FindStringSubmatch(nameWithOwner)
		owner = result[1]
		name = result[2]
	}

	return
}
