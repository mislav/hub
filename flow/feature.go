package flow

func FlowFeatureStart(featureName string) (err error) {
	branchName := "feature/" + featureName

	cmdGit := [][]string{}

	cmdGit1 := []string{"checkout", "develop"}
	cmdGit2 := []string{"checkout", "-b", branchName}

	cmdGit = append(cmdGit, cmdGit1, cmdGit2)

	err = launchCmdGit(cmdGit)

	return
}

func FlowFeatureFinish(featureName string) (err error) {
	branchName := "feature/" + featureName

	cmdGit := [][]string{}

	cmdGit1 := []string{"checkout", "develop"}
	cmdGit2 := []string{"merge", branchName, "--no-ff"}
	cmdGit3 := []string{"branch", "-d", branchName}

	cmdGit = append(cmdGit, cmdGit1, cmdGit2, cmdGit3)

	err = launchCmdGit(cmdGit)

	return
}

func FlowFeaturePullRequest(featureName string, params map[string][]string) (err error) {
	branchName := "feature/" + featureName
	messagePullRequest := "Pull request from " + branchName + " to develop"

	//hub pull-request -m 'test pull-request from cli' -b master -o boris
	cmdHub := []string{"pull-request", "-m", messagePullRequest, "-b", "develop"}

	if len(params["assignees"]) > 0 {
		cmdHub = append(cmdHub, "-a")
		for i := 0; i < len(params["assignees"]); i++ {
			cmdHub = append(cmdHub, params["assignees"][i])
		}
	}

	err = HubCmd(cmdHub...)

	return
}
