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
	Usage: "create [-p] [-d DESCRIPTION] [-h HOMEPAGE] [NAME]",
	Short: "Create this repository on GitHub and add GitHub as origin",
	Long: `Create a new public GitHub repository from the current git
repository and add remote origin at "git@github.com:USER/REPOSITORY.git";
USER is your GitHub username and REPOSITORY is the current working
directory name. To explicitly name the new repository, pass in NAME,
optionally in ORGANIZATION/NAME form to create under an organization
you're a member of. With -p, create a private repository, and with
-d and -h set the repository's description and homepage URL, respectively.
`,
}

var (
	flagCreatePrivate                         bool
	flagCreateDescription, flagCreateHomepage string
)

func init() {
	cmdCreate.Flag.BoolVarP(&flagCreatePrivate, "private", "p", false, "PRIVATE")
	cmdCreate.Flag.StringVarP(&flagCreateDescription, "description", "d", "", "DESCRIPTION")
	cmdCreate.Flag.StringVarP(&flagCreateHomepage, "homepage", "h", "", "HOMEPAGE")

	CmdRunner.Use(cmdCreate)
}

/*
  $ gh create
  ... create repo on github ...
  > git remote add -f origin git@github.com:YOUR_USER/CURRENT_REPO.git

  # with description:
  $ gh create -d 'It shall be mine, all mine!'

  $ gh create recipes
  [ repo created on GitHub ]
  > git remote add origin git@github.com:YOUR_USER/recipes.git

  $ gh create sinatra/recipes
  [ repo created in GitHub organization ]
  > git remote add origin git@github.com:sinatra/recipes.git
*/
func create(command *Command, args *Args) {
	_, err := git.Dir()
	if err != nil {
		err = fmt.Errorf("'create' must be run from inside a git repository")
		utils.Check(err)
	}

	var newRepoName string
	if args.IsParamsEmpty() {
		newRepoName, err = utils.DirName()
		utils.Check(err)
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

	var action string
	if gh.IsRepositoryExist(project) {
		ui.Printf("%s already exists on %s\n", project, project.Host)
		action = "set remote origin"
	} else {
		action = "created repository"
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
		args.Replace("git", "remote", "add", "-f", "origin", url)
	} else {
		args.Replace("git", "remote", "-v")
	}

	args.After("echo", fmt.Sprintf("%s:", action), project.String())
}
