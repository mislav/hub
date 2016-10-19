package flow

func getPrefixFeature() (prefix string, err error) {
	gitConfig, err := NewConfig()

	if err != nil {
		return
	}

	prefix, err = gitConfig.GetPrefix("feature")

	return
}

func getMergingBranchFeature() (branch string, err error) {
	gitConfig, err := NewConfig()

	if err != nil {
		return
	}

	branch, err = gitConfig.GetBranch("develop")

	return
}

func FlowFeatureStart(featureName string) (err error) {
	var prefixFeature string
	prefixFeature, err = getPrefixFeature()

	if err != nil {
		return
	}

	branchName := prefixFeature + featureName

	var mergingBranch string
	mergingBranch, err = getMergingBranchFeature()

	if err != nil {
		return
	}

	cmdGit := [][]string{}

	cmdGit1 := []string{"checkout", mergingBranch}
	cmdGit2 := []string{"checkout", "-b", branchName}

	cmdGit = append(cmdGit, cmdGit1, cmdGit2)

	err = launchCmdGit(cmdGit)

	return
}

func FlowFeatureFinish(featureName string) (err error) {
	var prefixFeature string
	prefixFeature, err = getPrefixFeature()

	if err != nil {
		return
	}

	branchName := prefixFeature + featureName

	var mergingBranch string
	mergingBranch, err = getMergingBranchFeature()

	if err != nil {
		return
	}

	cmdGit := [][]string{}

	cmdGit1 := []string{"checkout", mergingBranch}
	cmdGit2 := []string{"merge", branchName, "--no-ff"}
	cmdGit3 := []string{"branch", "-d", branchName}

	cmdGit = append(cmdGit, cmdGit1, cmdGit2, cmdGit3)

	err = launchCmdGit(cmdGit)

	return
}

func FlowFeaturePullRequest(featureName string, params map[string]string) (err error) {
	var prefixFeature string
	prefixFeature, err = getPrefixFeature()

	if err != nil {
		return
	}

	branchName := prefixFeature + featureName

	cmdGit := [][]string{}
	cmdGit1 := []string{"push", "origin", branchName}
	cmdGit = append(cmdGit, cmdGit1)

	err = launchCmdGit(cmdGit)

	if err != nil {
		return
	}

	messagePullRequest := "Pull request from " + branchName + " to develop"

	var mergingBranch string
	mergingBranch, err = getMergingBranchFeature()

	if err != nil {
		return
	}

	cmdHub := []string{"pull-request", "-m", messagePullRequest, "-b", mergingBranch}

	if len(params["assignees"]) > 0 {
		cmdHub = append(cmdHub, "-a", params["assignees"])
	}

	err = HubCmd(cmdHub...)

	return
}
