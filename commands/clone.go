package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var cmdClone = &Command{
	Run:          clone,
	GitExtension: true,
	Usage:        "clone [-p] [<OPTIONS>] [<USER>/]<REPOSITORY> [<DESTINATION>]",
	Long: `Clone a repository from GitHub.

## Options:
	-p
		(Deprecated) Clone private repositories over SSH.

	[<USER>/]<REPOSITORY>
		<USER> defaults to your own GitHub username.

	<DESTINATION>
		Directory name to clone into (default: <REPOSITORY>).

## Protocol used for cloning

The 'git:' protocol will be used for cloning public repositories, while the SSH
protocol will be used for private repositories and those that you have push
access to. Alternatively, hub can be configured to use HTTPS protocol for
everything. See "HTTPS instead of git protocol" and "HUB_PROTOCOL" of hub(1).

## Examples:
		$ hub clone rtomayko/ronn
		> git clone git://github.com/rtomayko/ronn.git

## See also:

hub-fork(1), hub(1), git-clone(1)
`,
}

func init() {
	CmdRunner.Use(cmdClone)
}

func clone(command *Command, args *Args) {
	if !args.IsParamsEmpty() {
		transformCloneArgs(args)
	}
}

func transformCloneArgs(args *Args) {
	isSSH := parseClonePrivateFlag(args)
	isCommandSubmodule := args.Command == "submodule"
	hasValueRegexp := regexp.MustCompile("^(--(upload-pack|template|depth|origin|branch|reference|name)|-[ubo])$")
	nameWithOwnerRegexp := regexp.MustCompile(NameWithOwnerRe)
	for i := 0; i < args.ParamsSize(); i++ {
		a := args.Params[i]

		if strings.HasPrefix(a, "-") {
			if hasValueRegexp.MatchString(a) {
				i++
			}
		} else {
			if nameWithOwnerRegexp.MatchString(a) && !isCloneable(a) {
				url := getProjectGitURL(a, isSSH, isCommandSubmodule)
				args.ReplaceParam(i, url)
			}

			break
		}
	}
}

func getProjectGitURL(nameWithOwner string, isSSH bool, isSubmodule bool) string {
	name, owner := parseCloneNameAndOwner(nameWithOwner)
	isMissingOwner := false
	host := getCloneHost()
	if owner == "" {
		isMissingOwner = true
		owner = host.User
	}

	var hostStr string
	if host != nil {
		hostStr = host.Host
	}

	expectWiki := strings.HasSuffix(name, ".wiki")
	if expectWiki {
		name = strings.TrimSuffix(name, ".wiki")
	}

	gh := github.NewClient(hostStr)

	var project *github.Project
	var repo *github.Repository
	if isMissingOwner {
		organizationsNames := parseUserOrganizationNames(gh)
		owners := append([]string{owner}, organizationsNames...)
		project, repo = determineRepository(owners, name, hostStr, gh)
	} else {
		project, repo = getRepository(owner, name, hostStr, gh, defaultRepositoryErrorHandler)
	}

	owner = repo.Owner.Login
	name = repo.Name
	if expectWiki {
		name = appendWikiToName(repo)
	}

	if !isSSH &&
		!isSubmodule &&
		!github.IsHttpsProtocol() {
		isSSH = repo.Private || repo.Permissions.Push
	}

	return project.GitURL(name, owner, isSSH)
}

func determineRepository(ownerCandidates []string, repositoryName string, hostString string, gh *github.Client) (project *github.Project, repo *github.Repository) {
	var handler errorHandler = func(err error, errorMessage string) {
		if !strings.Contains(err.Error(), "HTTP 404") {
			utils.Check(err)
		}
	}

	for _, owner := range ownerCandidates {
		project, repo = getRepository(owner, repositoryName, hostString, gh, handler)
		if repo != nil {
			break
		}
	}

	if repo == nil {
		err := fmt.Errorf("Error: repository doesn't exist for your username (%s) or any of the organizations you are part of (%s)", project.Owner, ownerCandidates)
		utils.Check(err)
	}

	return
}

func getRepository(owner string, repositoryName string, hostString string, gh *github.Client, handler errorHandler) (project *github.Project, repo *github.Repository) {
	if handler == nil {
		handler = defaultRepositoryErrorHandler
	}
	project = github.NewProject(owner, repositoryName, hostString)
	repo, err := gh.Repository(project)
	errorMessage := fmt.Sprintf("Error: repository %s/%s doesn't exist", project.Owner, project.Name)
	if err != nil {
		handler(err, errorMessage)
	}

	return
}

type errorHandler func(err error, errorMessage string)

func defaultRepositoryErrorHandler(err error, errorMessage string) {
	if strings.Contains(err.Error(), "HTTP 404") {
		err = fmt.Errorf(errorMessage)
		utils.Check(err)
	}
}

func appendWikiToName(repo *github.Repository) string {
	if !repo.HasWiki {
		utils.Check(fmt.Errorf("Error: %s/%s doesn't have a wiki", repo.Owner.Login, repo.Name))
	}
	return repo.Name + ".wiki"
}

func parseClonePrivateFlag(args *Args) bool {
	if i := args.IndexOfParam("-p"); i != -1 {
		args.RemoveParam(i)
		return true
	}

	return false
}

func parseCloneNameAndOwner(arg string) (name, owner string) {
	name, owner = arg, ""
	if strings.Contains(arg, "/") {
		split := strings.SplitN(arg, "/", 2)
		name = split[1]
		owner = split[0]
	}

	return
}

func getCloneHost() *github.Host {
	config := github.CurrentConfig()
	host, err := config.DefaultHost()
	if err != nil {
		utils.Check(github.FormatError("cloning repository", err))
	}

	return host
}

func parseUserOrganizationNames(gh *github.Client) (organizationsNames []string) {
	organizations, err := gh.FetchOrganizations()
	if err != nil {
		err = fmt.Errorf("Error: Problems fetching organizations for current user. %s", err)
		utils.Check(err)
	}

	for _, o := range organizations {
		organizationsNames = append(organizationsNames, o.Login)
	}

	return
}
