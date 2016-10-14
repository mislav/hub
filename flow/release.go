package flow

func FlowReleaseStart(featureName string) (err error) {
	branchName := "release/" + featureName

	cmdGit := [][]string{}

	cmdGit1 := []string{"checkout", "master"}
	cmdGit2 := []string{"checkout", "-b", branchName}

	cmdGit = append(cmdGit, cmdGit1, cmdGit2)

	err = launchCmdGit(cmdGit)

	return
}

func FlowReleaseFinish(featureName string) (err error) {
	branchName := "release/" + featureName

	cmdGit := [][]string{}

	cmdGit1 := []string{"checkout", "master"}
	cmdGit2 := []string{"merge", branchName}
	cmdGit3 := []string{"tag", "-a", featureName}
	cmdGit4 := []string{"checkout", "develop"}
	cmdGit5 := []string{"merge", branchName}
	cmdGit6 := []string{"checkout", "master"}

	cmdGit = append(cmdGit, cmdGit1, cmdGit2, cmdGit3, cmdGit4, cmdGit5, cmdGit6)

	err = launchCmdGit(cmdGit)

	return
}
