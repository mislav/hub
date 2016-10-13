package flow

import "github.com/github/hub/git"

func flowFeatureStart(featureName string) (err error) {
	branchName := "feature/" + featureName

	cmdGit := [][]string{}

	cmdGit1 := []string{"checkout", "develop"}
	cmdGit2 := []string{"checkout", "-b", branchName}

	cmdGit = append(cmdGit, cmdGit1, cmdGit2)

	err = launchCmdGit(cmdGit)

	return
}

func flowFeatureFinish(featureName string) (err error) {
	branchName := "feature/" + featureName

	cmdGit := [][]string{}

	cmdGit1 := []string{"checkout", "develop"}
	cmdGit2 := []string{"merge", branchName}
	cmdGit3 := []string{"branch", "-d", branchName}

	cmdGit = append(cmdGit, cmdGit1, cmdGit2, cmdGit3)

	err = launchCmdGit(cmdGit)

	return
}

func launchCmdGit(cmdGit [][]string) (err error) {
	for i := range cmdGit {
		err = git.Spawn(cmdGit[i]...)

		if err != nil {
			break
		}
	}

	return
}
