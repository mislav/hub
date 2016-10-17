package flow

import (
	"github.com/github/hub/cmd"
	"github.com/github/hub/git"
)

func launchCmdGit(cmdGit [][]string) (err error) {
	for i := range cmdGit {
		err = git.Spawn(cmdGit[i]...)

		if err != nil {
			break
		}
	}

	return
}

func HubCmd(args ...string) (err error) {
	cmd := cmd.New("hub")

	for _, a := range args {
		cmd.WithArg(a)
	}

	return cmd.Spawn()
}
