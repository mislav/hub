package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/github/hub/v2/git"
	"github.com/github/hub/v2/github"
	"github.com/github/hub/v2/utils"
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
		(Deprecated) Use the ''ssh:'' protocol instead of ''git:'' for the remote URL.
		The writeable ''ssh:'' protocol is automatically used for own repos, GitHub
		Enterprise remotes, and private or pushable repositories.

	<USER>[/<REPOSITORY>]
		If <USER> is "origin", that value will be substituted for your GitHub
		username. <REPOSITORY> defaults to the name of the current working directory.

## Examples:
		$ hub remote add jingweno
		> git remote add jingweno git://github.com/jingweno/REPO.git

		$ hub remote add origin
		> git remote add origin git@github.com:USER/REPO.git

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

	re := regexp.MustCompile(fmt.Sprintf(`^%s(/%s)?$`, OwnerRe, NameRe))
	if !re.MatchString(ownerWithName) {
		return
	}
	owner := ownerWithName
	name := ""
	if strings.Contains(ownerWithName, "/") {
		parts := strings.SplitN(ownerWithName, "/", 2)
		owner, name = parts[0], parts[1]
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

	p := utils.NewArgsParser()
	p.RegisterValue("-t")
	p.RegisterValue("-m")
	params, _ := p.Parse(args.Params)
	if len(params) > 3 {
		return
	}

	for i, pi := range p.PositionalIndices {
		if i == 1 && strings.Contains(params[i], "/") {
			args.ReplaceParam(pi, owner)
		} else if i == 2 {
			args.RemoveParam(pi)
		}
	}
	if len(params) == 2 && owner == "origin" {
		owner = hostConfig.User
	}

	if strings.EqualFold(owner, hostConfig.User) {
		owner = hostConfig.User
	}

	project := github.NewProject(owner, name, host)

	isPrivate := parseRemotePrivateFlag(args) || owner == hostConfig.User || project.Host != github.GitHubHost
	if !isPrivate {
		gh := github.NewClient(project.Host)
		repo, err := gh.Repository(project)
		if err != nil {
			if strings.Contains(err.Error(), "HTTP 404") {
				err = fmt.Errorf("Error: repository %s/%s doesn't exist", project.Owner, project.Name)
			}
			utils.Check(err)
		}
		isPrivate = repo.Private || repo.Permissions.Push
	}

	url := project.GitURL("", "", isPrivate)
	args.AppendParams(url)
}

func parseRemotePrivateFlag(args *Args) bool {
	if i := args.IndexOfParam("-p"); i != -1 {
		args.RemoveParam(i)
		return true
	}

	return false
}
