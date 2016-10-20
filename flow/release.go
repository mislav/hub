package flow

func getPrefixRelease() (prefix string, err error) {
	gitConfig, err := NewConfig()

	if err != nil {
		return
	}

	prefix, err = gitConfig.GetPrefix("release")

	return
}

func getStartBranchRelease() (branch string, err error) {
	gitConfig, err := NewConfig()

	if err != nil {
		return
	}

	branch, err = gitConfig.GetBranch("develop")

	return
}

func getFinishBranchRelease() (branchMaster string, branchDevelop string, err error) {
	gitConfig, err := NewConfig()

	if err != nil {
		return
	}

	branchDevelop, err = gitConfig.GetBranch("develop")
	branchMaster, err = gitConfig.GetBranch("master")

	return
}

func FlowReleaseStart(releaseName string) (err error) {
	var prefixRelease string
	prefixRelease, err = getPrefixRelease()

	if err != nil {
		return
	}

	branchName := prefixRelease + releaseName

	var startingBranch string
	startingBranch, err = getStartBranchRelease()

	if err != nil {
		return
	}

	cmdGit := [][]string{}

	cmdGit1 := []string{"checkout", startingBranch}
	cmdGit2 := []string{"checkout", "-b", branchName}

	cmdGit = append(cmdGit, cmdGit1, cmdGit2)

	err = launchCmdGit(cmdGit)

	return
}

func FlowReleaseFinish(releaseName string) (err error) {
	var prefixRelease string
	prefixRelease, err = getPrefixRelease()

	if err != nil {
		return
	}

	branchName := prefixRelease + releaseName

	var (
		masterBranch,
		developBranch string
	)

	masterBranch, developBranch, err = getFinishBranchRelease()

	if err != nil {
		return
	}

	cmdGit := [][]string{}

	cmdGit1 := []string{"checkout", masterBranch}
	cmdGit2 := []string{"merge", branchName, "--no-ff"}
	cmdGit3 := []string{"tag", "-a", releaseName}
	cmdGit4 := []string{"checkout", developBranch}
	cmdGit5 := []string{"merge", branchName, "--no-ff"}
	cmdGit6 := []string{"checkout", masterBranch}
	cmdGit7 := []string{"branch", "-d", branchName}

	cmdGit = append(cmdGit, cmdGit1, cmdGit2, cmdGit3, cmdGit4, cmdGit5, cmdGit6, cmdGit7)

	err = launchCmdGit(cmdGit)

	return
}

func FlowReleasePullRequest(releaseName string, params map[string]string) (err error) {
	var prefixRelease string
	prefixRelease, err = getPrefixRelease()

	if err != nil {
		return
	}

	branchName := prefixRelease + releaseName

	cmdGit := [][]string{}
	cmdGit1 := []string{"push", "origin", branchName}
	cmdGit = append(cmdGit, cmdGit1)

	err = launchCmdGit(cmdGit)

	if err != nil {
		return
	}

	var (
		masterBranch,
		developBranch string
	)

	masterBranch, developBranch, err = getFinishBranchRelease()

	if err != nil {
		return
	}

	messagePullRequestMaster := "Pull request from " + branchName + " to " + masterBranch
	messagePullRequestDevelop := "Pull request from " + branchName + " to " + developBranch

	cmdHub := [][]string{}

	cmdHub1 := []string{"pull-request", "-m", messagePullRequestMaster, "-b", masterBranch}
	cmdHub2 := []string{"pull-request", "-m", messagePullRequestDevelop, "-b", developBranch}

	cmdHub = append(cmdHub, cmdHub1, cmdHub2)

	if len(params["assignees"]) > 0 {
		for _, cmd := range cmdHub {
			cmd = append(cmd, "-a", params["assignees"])
		}
	}

	err = launchCmdHub(cmdHub)

	return
}
