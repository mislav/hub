package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/github/hub/git"
	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var cmdRemote = &Command{
	Run:          remote,
	GitExtension: true,
	Usage: `
remote add [-p] [<OPTIONS>] <USER>[/<REPOSITORY>]
remote set-url [-p] [<OPTIONS>] <NAME> <USER>[/<REPOSITORY>]
`,
	Long: `Add a git remote for a GitHub repository.

## Options:
	-p
		(Deprecated) Use the 'ssh:' protocol for the remote URL.

	<USER>[/<REPOSITORY>]
		If <USER> is "origin", that value will be substituted for your GitHub
		username. <REPOSITORY> defaults to the name of the current working directory.

## Examples:
		$ hub remote add jingweno
		> git remote add jingweno git://github.com/jingweno/REPO.git

		$ hub remote add origin
		> git remote add origin git://github.com/USER/REPO.git

## See also:

hub-fork(1), hub(1), git-remote(1)
`,
}

func init() {
	CmdRunner.Use(cmdRemote)
}

/*
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

	var host string
	mainProject, err := localRepo.MainProject()
	if err == nil {
		host = mainProject.Host
	}

	if name == "" {
		if mainProject != nil {
			name = mainProject.Name
		} else {
			dirName, err := git.WorkdirName()
			utils.Check(err)
			name = github.SanitizeProjectName(dirName)
		}
	}

	var hostConfig *github.Host
	if host == "" {
		hostConfig, err = github.CurrentConfig().DefaultHost()
	} else {
		hostConfig, err = github.CurrentConfig().PromptForHost(host)
	}
	if err != nil {
		utils.Check(github.FormatError("adding remote", err))
	}
	host = hostConfig.Host

	words := args.Words()
	isPrivate := parseRemotePrivateFlag(args)
	if len(words) == 2 && words[1] == "origin" {
		// Origin special case triggers default user/repo
		owner = hostConfig.User
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
