package commands

import (
	"fmt"
	"os"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/utils"
)

var cmdRemoteAdd = &Command{
	Run: remoteAdd,
	Usage: "remote add [-p] USER",
	Short: "Add remote from GitHub repository",
	Long: `Add remote from GitHub repository, using USER as the username and the current repository name.
If -p is provided, the SSH remote will be added.
If USER is "origin", your own username will be used.
`,
}

func toSSHOrNotToSSH(args []string) bool {
	for i:=0;i<len(args);i++ {
		if args[i] == "-p" {
			return true
		}
	}
	return false
}

func remoteAdd(command *Command, args []string) {
	if len(args) == 0 || len(args) > 0 && args[0] != "add" || len(args) == 2 && args[1] == "-p" || len(args) > 3 {
		command.PrintUsage()
		os.Exit(1)
	}

	var name string

	if args[1] == "-p" {
		name = args[2]
	} else {
		name = args[1]
	}
	
	flagRemoteAddSSH := toSSHOrNotToSSH(args)

	gh := github.New()
	project := gh.Project

	if name == "origin" {
		project.Owner = gh.FetchUsername()
	} else {
		project.Owner = name
	}

	var url string

	if flagRemoteAddSSH {
		url = project.SshURL()
	} else {
		url = project.GitURL()
	}

	err := git.AddRemote(name, url)
	utils.Check(err)
	fmt.Printf("The remote %s has been added.\n", url)
}
