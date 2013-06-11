package commands

import (
  "github.com/jingweno/gh/github"
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

  err := gh.ForkRepository(project.Name, project.Owner)


}
