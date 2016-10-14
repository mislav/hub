package flow

import "github.com/github/hub/git"

func launchCmdGit(cmdGit [][]string) (err error) {
	for i := range cmdGit {
		err = git.Spawn(cmdGit[i]...)

		if err != nil {
			break
		}
	}

	return
}
