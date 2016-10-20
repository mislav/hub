package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/github/hub/git"
	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var cmdCreate = &Command{
	Run:   create,
	Usage: "create [-poc] [-d <DESCRIPTION>] [-h <HOMEPAGE>] [[<ORGANIZATION>/]<NAME>]",
	Long: `Create a new repository on GitHub and add a git remote for it.

## Options:
	-p
		Create a private repository.

	-d <DESCRIPTION>
		Use this text as the description of the GitHub repository.

	-h <HOMEPAGE>
		Use this text as the URL of the GitHub repository.

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

var (
	flagCreatePrivate,
	flagCreateBrowse,
	flagCreateCopy bool

	flagCreateDescription,
	flagCreateHomepage string
)

func init() {
	cmdCreate.Flag.BoolVarP(&flagCreatePrivate, "private", "p", false, "PRIVATE")
	cmdCreate.Flag.BoolVarP(&flagCreateBrowse, "browse", "o", false, "BROWSE")
	cmdCreate.Flag.BoolVarP(&flagCreateCopy, "copy", "c", false, "COPY")
	cmdCreate.Flag.StringVarP(&flagCreateDescription, "description", "d", "", "DESCRIPTION")
	cmdCreate.Flag.StringVarP(&flagCreateHomepage, "homepage", "h", "", "HOMEPAGE")

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
		reg := regexp.MustCompile("^[^-]")
		if !reg.MatchString(args.FirstParam()) {
			err = fmt.Errorf("invalid argument: %s", args.FirstParam())
			utils.Check(err)
		}
		newRepoName = args.FirstParam()
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

	if gh.IsRepositoryExist(project) {
		ui.Errorln("Existing repository detected. Updating git remote")
	} else {
		if !args.Noop {
			repo, err := gh.CreateRepository(project, flagCreateDescription, flagCreateHomepage, flagCreatePrivate)
			utils.Check(err)
			project = github.NewProject(repo.FullName, "", project.Host)
		}
	}

	localRepo, err := github.LocalRepo()
	utils.Check(err)

	remote, _ := localRepo.OriginRemote()
	if remote == nil || remote.Name != "origin" {
		url := project.GitURL("", "", true)
		args.Before("git", "remote", "add", "-f", "origin", url)
	}

	webUrl := project.WebURL("", "", "")
	args.NoForward()
	printBrowseOrCopy(args, webUrl, flagCreateBrowse, flagCreateCopy)
}
