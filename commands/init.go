package commands

import (
	"path/filepath"
	"regexp"
	"strings"

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
	err := transformInitArgs(args)
	utils.Check(err)
}

func transformInitArgs(args *Args) error {
	if !parseInitFlag(args) {
		return nil
	}

	var err error
	dirToInit := "."
	hasValueRegxp := regexp.MustCompile("^--(template|separate-git-dir|shared)$")

	// Find the first argument that isn't related to any of the init flags.
	// We assume this is the optional `directory` argument to git init.
	for i := 0; i < args.ParamsSize(); i++ {
		arg := args.Params[i]
		if hasValueRegxp.MatchString(arg) {
			i++
		} else if !strings.HasPrefix(arg, "-") {
			dirToInit = arg
			break
		}
	}

	dirToInit, err = filepath.Abs(dirToInit)
	if err != nil {
		return err
	}

	// Assume that the name of the working directory is going to be the name of
	// the project on GitHub.
	projectName := strings.Replace(filepath.Base(dirToInit), " ", "-", -1)
	project := github.NewProject("", projectName, "")
	url := project.GitURL("", "", true)

	addRemote := []string{
		"git", "--git-dir", filepath.Join(dirToInit, ".git"),
		"remote", "add", "origin", url,
	}
	args.After(addRemote...)

	return nil
}

func parseInitFlag(args *Args) bool {
	if i := args.IndexOfParam("-g"); i != -1 {
		args.RemoveParam(i)
		return true
	}

	return false
}
