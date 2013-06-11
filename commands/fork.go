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

func fork(cmd *Command, args []string) {
	gh := github.New()
	project := gh.Project

	newRemote, err := gh.ForkRepository(project.Name, project.Owner, flagForkNoRemote)
	utils.Check(err)

	if !flagForkNoRemote && newRemote != "" {
		fmt.Printf("New remote: %s\n", newRemote)
	}
}
