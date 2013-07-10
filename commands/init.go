package commands

import (
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
)

var cmdInit = &Command{
	Run:          gitInit,
	GitExtension: true,
	Usage:        "init -g",
	Short:        "Create an empty git repository or reinitialize an existing one",
	Long: `Create a git repository as with git-init(1) and add remote origin at
"git@github.com:USER/REPOSITORY.git"; USER is your GitHub username and
REPOSITORY is the current working directory's basename.
`,
}

/*
  $ gh init -g
  > git init
  > git remote add origin git@github.com:USER/REPO.git
*/
func gitInit(command *Command, args *Args) {
	if !args.IsParamsEmpty() {
		err := transformInitArgs(args)
		utils.Check(err)
	}
}

func transformInitArgs(args *Args) error {
	if !parseInitFlag(args) {
		return nil
	}

	name, err := utils.DirName()
	if err != nil {
		return err
	}

	owner := github.CurrentConfig().FetchUser()
	project := github.Project{Owner: owner, Name: name}
	url := project.GitURL(name, owner, true)
	args.After("git", "remote", "add", "origin", url)

	return nil
}

func parseInitFlag(args *Args) bool {
	if i := args.IndexOfParam("-g"); i != -1 {
		args.RemoveParam(i)
		return true
	}

	return false
}
