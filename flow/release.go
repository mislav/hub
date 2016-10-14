package flow

func FlowReleaseStart(releaseName string) (err error) {
	branchName := "release/" + releaseName

	cmdGit := [][]string{}

	cmdGit1 := []string{"checkout", "master"}
	cmdGit2 := []string{"checkout", "-b", branchName}

	cmdGit = append(cmdGit, cmdGit1, cmdGit2)

	err = launchCmdGit(cmdGit)

	return
}

func FlowReleaseFinish(releaseName string) (err error) {
	branchName := "release/" + releaseName

	cmdGit := [][]string{}

	cmdGit1 := []string{"checkout", "master"}
	cmdGit2 := []string{"merge", "--no-ff", branchName}
	cmdGit3 := []string{"tag", "-a", releaseName}
	cmdGit4 := []string{"checkout", "develop"}
	cmdGit5 := []string{"merge", "--no-ff", branchName}
	cmdGit6 := []string{"checkout", "master"}

	cmdGit = append(cmdGit, cmdGit1, cmdGit2, cmdGit3, cmdGit4, cmdGit5, cmdGit6)

	err = launchCmdGit(cmdGit)

	return
}
