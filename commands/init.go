package commands

import (
	"path/filepath"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
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

func init() {
	CmdRunner.Use(cmdInit)
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

	var (
		name   string
		newDir bool
		err    error
	)

	if args.IsParamsEmpty() {
		name, err = utils.DirName()
		if err != nil {
			return err
		}
	} else {
		name = args.LastParam()
		newDir = true
	}

	project := github.NewProject("", name, "")
	url := project.GitURL("", "", true)

	cmds := []string{"git"}
	if newDir {
		cmds = append(cmds, "--git-dir", filepath.Join(name, ".git"))
	}

	cmds = append(cmds, "remote", "add", "origin", url)
	args.After(cmds...)

	return nil
}

func parseInitFlag(args *Args) bool {
	if i := args.IndexOfParam("-g"); i != -1 {
		args.RemoveParam(i)
		return true
	}

	return false
}
