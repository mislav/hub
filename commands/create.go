package commands

import (
	"fmt"
	"strings"

	"github.com/github/hub/v2/git"
	"github.com/github/hub/v2/github"
	"github.com/github/hub/v2/ui"
	"github.com/github/hub/v2/utils"
)

var cmdCreate = &Command{
	Run:   create,
	Usage: "create [-poc] [-d <DESCRIPTION>] [-h <HOMEPAGE>] [[<ORGANIZATION>/]<NAME>]",
	Long: `Create a new repository on GitHub and add a git remote for it.

## Options:
	-p, --private
		Create a private repository.

	-d, --description <DESCRIPTION>
		A short description of the GitHub repository.

	-h, --homepage <HOMEPAGE>
		A URL with more information about the repository. Use this, for example, if
		your project has an external website.

	--remote-name <REMOTE>
		Set the name for the new git remote (default: "origin").

	-o, --browse
		Open the new repository in a web browser.

	-c, --copy
		Put the URL of the new repository to clipboard instead of printing it.

	[<ORGANIZATION>/]<NAME>
		The name for the repository on GitHub (default: name of the current working
		directory).

		Optionally, create the repository within <ORGANIZATION>.

## Examples:
		$ hub create
		[ repo created on GitHub ]
		> git remote add -f origin git@github.com:USER/REPO.git

		$ hub create sinatra/recipes
		[ repo created in GitHub organization ]
		> git remote add -f origin git@github.com:sinatra/recipes.git

## See also:

hub-init(1), hub(1)
`,
}

func init() {
	CmdRunner.Use(cmdCreate)
}

func create(command *Command, args *Args) {
	_, err := git.Dir()
	if err != nil {
		err = fmt.Errorf("'create' must be run from inside a git repository")
		utils.Check(err)
	}

	var newRepoName string
	if args.IsParamsEmpty() {
		dirName, err := git.WorkdirName()
		utils.Check(err)
		newRepoName = github.SanitizeProjectName(dirName)
	} else {
		newRepoName = args.FirstParam()
		if newRepoName == "" {
			utils.Check(command.UsageError(""))
		}
	}

	config := github.CurrentConfig()
	host, err := config.DefaultHost()
	if err != nil {
		utils.Check(github.FormatError("creating repository", err))
	}

	owner := host.User
	if strings.Contains(newRepoName, "/") {
		split := strings.SplitN(newRepoName, "/", 2)
		owner = split[0]
		newRepoName = split[1]
	}

	project := github.NewProject(owner, newRepoName, host.Host)
	gh := github.NewClient(project.Host)

	flagCreatePrivate := args.Flag.Bool("--private")

	repo, err := gh.Repository(project)
	if err == nil {
		foundProject := github.NewProject(repo.FullName, "", project.Host)
		if foundProject.SameAs(project) {
			if !repo.Private && flagCreatePrivate {
				err = fmt.Errorf("Repository '%s' already exists and is public", repo.FullName)
				utils.Check(err)
			} else {
				ui.Errorln("Existing repository detected")
				project = foundProject
			}
		} else {
			repo = nil
		}
	} else {
		repo = nil
	}

	if repo == nil {
		if !args.Noop {
			flagCreateDescription := args.Flag.Value("--description")
			flagCreateHomepage := args.Flag.Value("--homepage")
			repo, err := gh.CreateRepository(project, flagCreateDescription, flagCreateHomepage, flagCreatePrivate)
			utils.Check(err)
			project = github.NewProject(repo.FullName, "", project.Host)
		}
	}

	localRepo, err := github.LocalRepo()
	utils.Check(err)

	originName := args.Flag.Value("--remote-name")
	if originName == "" {
		originName = "origin"
	}

	if originRemote, err := localRepo.RemoteByName(originName); err == nil {
		originProject, err := originRemote.Project()
		if err != nil || !originProject.SameAs(project) {
			ui.Errorf("A git remote named '%s' already exists and is set to push to '%s'.\n", originRemote.Name, originRemote.PushURL)
		}
	} else {
		url := project.GitURL("", "", true)
		args.Before("git", "remote", "add", "-f", originName, url)
	}

	webURL := project.WebURL("", "", "")
	args.NoForward()
	flagCreateBrowse := args.Flag.Bool("--browse")
	flagCreateCopy := args.Flag.Bool("--copy")
	printBrowseOrCopy(args, webURL, flagCreateBrowse, flagCreateCopy)
}
