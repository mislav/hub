package commands

import (
	"fmt"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"os"
)

var cmdRemote = &Command{
	Run:   remote,
	Usage: "remote [-p] OPTIONS USER[/REPOSITORY]",
	Short: "View and manage a set of remote repositories",
	Long: `Add remote "git://github.com/USER/REPOSITORY.git" as with
git-remote(1). When /REPOSITORY is omitted, the basename of the
current working directory is used. With -p, use private remote
"git@github.com:USER/REPOSITORY.git". If USER is "origin"
then uses your GitHub login.
`,
}

var flagRemoteAddSSH bool

func init() {
	cmdRemote.Flag.BoolVar(&flagRemoteAddSSH, "p", false, "")
}

func remote(command *Command, args []string) {
	if len(args) <= 1 || len(args) >= 3 || len(args) > 0 && args[0] != "add" {
		command.PrintUsage()
		os.Exit(1)
	}

	name := args[1]

	gh := github.New()
	url, err := gh.RemoteAdd(name, flagRemoteAddSSH)
	utils.Check(err)
	fmt.Printf("The remote %s has been added.\n", url)
}
