package commands

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/github/hub/v2/github"
	"github.com/github/hub/v2/utils"
)

var cmdInit = &Command{
	Run:          gitInit,
	GitExtension: true,
	Usage:        "init -g",
	Long: `Initialize a git repository and add a remote pointing to GitHub.

## Options:
	-g
		After initializing the repository locally, add the "origin" remote pointing
		to "<USER>/<REPO>" repository on GitHub.

		<USER> is your GitHub username, while <REPO> is the name of the current
		working directory.

## Examples:
		$ hub init -g
		> git init
		> git remote add origin git@github.com:USER/REPO.git

## See also:

hub-create(1), hub(1), git-init(1)
`,
}

func init() {
	CmdRunner.Use(cmdInit)
}

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
	hasValueRegexp := regexp.MustCompile("^--(template|separate-git-dir|shared)$")

	// Find the first argument that isn't related to any of the init flags.
	// We assume this is the optional `directory` argument to git init.
	for i := 0; i < args.ParamsSize(); i++ {
		arg := args.Params[i]
		if hasValueRegexp.MatchString(arg) {
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

	config := github.CurrentConfig()
	host, err := config.DefaultHost()
	if err != nil {
		utils.Check(github.FormatError("initializing repository", err))
	}

	// Assume that the name of the working directory is going to be the name of
	// the project on GitHub.
	projectName := strings.Replace(filepath.Base(dirToInit), " ", "-", -1)
	project := github.NewProject(host.User, projectName, host.Host)
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
