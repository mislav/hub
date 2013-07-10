package commands

import (
	"fmt"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
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
	cmdCreate.Flag.BoolVar(&flagCreatePrivate, "p", false, "PRIVATE")
	cmdCreate.Flag.StringVar(&flagCreateDescription, "d", "", "DESCRIPTION")
	cmdCreate.Flag.StringVar(&flagCreateHomepage, "h", "", "HOMEPAGE")
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
	var (
		name string
		err  error
	)
	if args.IsParamsEmpty() {
		name, err = utils.DirName()
		utils.Check(err)
	} else {
		name = args.FirstParam()
	}

	var msg string
	project := github.NewProjectFromNameAndOwner(name, "")
	gh := github.NewWithoutProject()
	if gh.IsRepositoryExist(project) {
		fmt.Printf("%s already exists on %s\n", project, github.GitHubHost)
		msg = "set remmote origin"
	} else {
		if !args.Noop {
			repo, err := gh.CreateRepository(project, flagCreateDescription, flagCreateHomepage, flagCreatePrivate)
			utils.Check(err)
			project = github.NewProjectFromNameAndOwner("", repo.FullName)
		}
		msg = "created repository"
	}

	remote, _ := git.OriginRemote()
	if remote == nil {
		url := project.GitURL("", "", true)
		args.Replace("git", "remote", "add", "-f", "origin", url)
	} else {
		args.Replace("git", "remote", "-v")
	}

	args.After("echo", fmt.Sprintf("%s:", msg), project.String())
}
