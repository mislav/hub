package commands

import (
	"fmt"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"strings"
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
	gh := github.NewWithoutProject()
	var nameWithOwner string
	if args.IsParamsEmpty() {
		name, err := repoName()
		utils.Check(err)
		nameWithOwner = fmt.Sprintf("%s/%s", gh.Config.FetchUser(), name)
	} else {
		nameWithOwner = args.FirstParam()
	}

	owner, name := parseCreateOwnerAndName(nameWithOwner)
	if owner == "" {
		owner = gh.Config.FetchUser()
	}
	project := github.Project{Name: name, Owner: owner}

	var msg string
	if gh.IsRepositoryExist(project) {
		fmt.Printf("%s already exists on %s\n", project, github.GitHubHost)
		msg = "set remmote origin"
	} else {
		msg = "created repository"
		if !args.Noop {
			repo, err := gh.CreateRepository(project, flagCreateDescription, flagCreateHomepage, flagCreatePrivate)
			utils.Check(err)
			owner, name = parseCreateOwnerAndName(repo.FullName)
			project = github.Project{Name: name, Owner: owner}
		}
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

func parseCreateOwnerAndName(nameWithOwner string) (owner, name string) {
	if strings.Contains(nameWithOwner, "/") {
		result := strings.SplitN(nameWithOwner, "/", 2)
		owner = result[0]
		name = result[1]
	} else {
		name = nameWithOwner
	}

	return
}
