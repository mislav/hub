package commands

import (
	"fmt"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
)

var cmdFork = &Command{
	Run:   fork,
	Usage: "fork [--no-remote]",
	Short: "Make a fork of a remote repository on GitHub and add as remote",
	Long: `Forks the original project (referenced by "origin" remote) on GitHub and
adds a new remote for it under your username.
`,
}

var flagForkNoRemote bool

func init() {
	cmdFork.Flag.BoolVar(&flagForkNoRemote, "no-remote", false, "")
}

/*
  $ gh fork
  [ repo forked on GitHub ]
  > git remote add -f YOUR_USER git@github.com:YOUR_USER/CURRENT_REPO.git

  $ gh fork --no-remote
  [ repo forked on GitHub ]
*/
func fork(cmd *Command, args *Args) {
	gh := github.New()
	project := gh.Project

	var forkURL string
	if args.Noop {
		args.Before(fmt.Sprintf("Would request a fork to %s:%s", project.Owner, project.Name), "")
		forkURL = "FORK_URL"
	} else {
		repo, err := gh.ForkRepository(project.Name, project.Owner, flagForkNoRemote)
		utils.Check(err)

		forkURL = repo.SshURL
	}

	if !flagForkNoRemote {
		newRemote := gh.Config.User
		args.Replace("git", "remote", "add", "-f", newRemote, forkURL)
		args.After("echo", fmt.Sprintf("New remote: %s", newRemote))
	}
}
