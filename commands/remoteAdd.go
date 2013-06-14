package commands

import (
	"fmt"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"os"
)

var cmdRemoteAdd = &Command{
	Run:   remoteAdd,
	Usage: "remote [-p] add USER",
	Short: "Add remote from GitHub repository",
	Long: `Add remote from GitHub repository, using USER as the username and the current repository name.
If -p is provided, the SSH remote will be added.
If USER is "origin", your own username will be used.
`,
}

var flagRemoteAddSSH bool

func init() {
	cmdRemoteAdd.Flag.BoolVar(&flagRemoteAddSSH, "p", false, "")
}

func remoteAdd(command *Command, args []string) {
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
