package commands

import (
	"fmt"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"regexp"
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

	configs := github.CurrentConfigs()
	credentials := configs.DefaultCredentials()

	owner := credentials.User
	if strings.Contains(newRepoName, "/") {
		split := strings.SplitN(newRepoName, "/", 2)
		owner = split[0]
		newRepoName = split[1]
	}

	project := github.NewProject(owner, newRepoName, credentials.Host)
	gh := github.NewClient(project)

	var action string
	if gh.IsRepositoryExist(project) {
		fmt.Printf("%s already exists on %s\n", project, project.Host)
		action = "set remote origin"
	} else {
		action = "created repository"
		if !args.Noop {
			repo, err := gh.CreateRepository(project, flagCreateDescription, flagCreateHomepage, flagCreatePrivate)
			utils.Check(err)
			project = github.NewProject(repo.FullName, "", project.Host)
		}
	}

	remote, _ := git.OriginRemote()
	if remote == nil {
		url := project.GitURL("", "", true)
		args.Replace("git", "remote", "add", "-f", "origin", url)
	} else {
		args.Replace("git", "remote", "-v")
	}

	args.After("echo", fmt.Sprintf("%s:", action), project.String())
}
