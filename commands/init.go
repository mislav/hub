package commands

import (
	"path/filepath"

	flag "github.com/github/hub/Godeps/_workspace/src/github.com/ogier/pflag"
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
	f := parseInitFlag(args.Params)
	if !f.AddRemote {
		return nil
	}

	// remove "-g" from args so that `git-init` wont complain
	if i := args.IndexOfParam("-g"); i != -1 {
		args.RemoveParam(i)
	}

	var (
		dir    string
		newDir bool
		err    error
	)

	if f.Dir == "" {
		dir, err = utils.DirName()
		if err != nil {
			return err
		}
	} else {
		dir = f.Dir
		newDir = true
	}

	project := github.NewProject("", dir, "")
	url := project.GitURL("", "", true)

	cmds := []string{"git"}
	if newDir {
		cmds = append(cmds, "--git-dir", filepath.Join(dir, ".git"))
	}

	cmds = append(cmds, "remote", "add", "origin", url)
	args.After(cmds...)

	return nil
}

type initFlag struct {
	AddRemote bool
	Dir       string
}

func parseInitFlag(params []string) *initFlag {
	var (
		initFlagSet flag.FlagSet

		initFlag       = &initFlag{}
		quiet          bool
		bare           bool
		template       string
		separateGitDir string
		shared         string
	)

	initFlagSet.BoolVarP(&initFlag.AddRemote, "", "g", false, "")
	initFlagSet.BoolVarP(&quiet, "quiet", "q", false, "")
	initFlagSet.BoolVar(&bare, "bare", false, "")
	initFlagSet.StringVar(&template, "template", "", "")
	initFlagSet.StringVar(&separateGitDir, "separate-git-dir", "", "")
	initFlagSet.StringVar(&shared, "shared", "", "")

	err := initFlagSet.Parse(params)
	utils.Check(err)

	a := initFlagSet.Args()
	if len(a) != 0 {
		initFlag.Dir = a[0]
	}

	return initFlag
}
