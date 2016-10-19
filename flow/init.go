package flow

import (
	"fmt"

	"github.com/github/hub/git"
	"github.com/github/hub/utils"
)

type BranchingFlow struct {
	master        string
	develop       string
	featurePrefix string
	releasePrefix string
	hotfixPrefix  string
}

func getBranching() (branching BranchingFlow) {
	fmt.Print("Which branch for master [master]: ")
	branching.master = "master"
	fmt.Scanln(&branching.master)

	fmt.Print("Which branch for develop [develop]: ")
	branching.develop = "develop"
	fmt.Scanln(&branching.develop)

	fmt.Println()

	fmt.Print("Which prefix for feature branch [feature/]: ")
	branching.featurePrefix = "feature/"
	fmt.Scanln(&branching.featurePrefix)

	fmt.Print("Which prefix for release branch [release/]: ")
	branching.releasePrefix = "release/"
	fmt.Scanln(&branching.releasePrefix)

	fmt.Print("Which prefix for hotfix branch [hotfix/]: ")
	branching.hotfixPrefix = "hotfix/"
	fmt.Scanln(&branching.hotfixPrefix)

	return
}

func updateConfigFile(branching BranchingFlow) (err error) {
	gitConfig, err := NewConfig()

	if err != nil {
		utils.Check(fmt.Errorf("%s", err.Error()))
	}

	err = gitConfig.CreateSections(gitConfig.sectionBranch, gitConfig.sectionPrefix)

	if err != nil {
		return
	}

	gitConfig.CreateKey(gitConfig.sectionBranch, "master", branching.master)
	gitConfig.CreateKey(gitConfig.sectionBranch, "develop", branching.develop)

	gitConfig.CreateKey(gitConfig.sectionPrefix, "feature", branching.featurePrefix)
	gitConfig.CreateKey(gitConfig.sectionPrefix, "release", branching.releasePrefix)
	gitConfig.CreateKey(gitConfig.sectionPrefix, "hotfix", branching.hotfixPrefix)

	gitConfig.Save()

	return
}

func initMasterBranch() (err error) {
	git.Quiet("checkout", "master")

	_, err = git.Ref("HEAD")

	if err != nil {
		git.Run("commit", "-m", "Initial commit", "--allow-empty")
	}

	return
}

func initDevelopBranch() (err error) {
	var branches []string

	branches, err = git.LocalBranches()

	if err != nil {
		return
	}

	developExist := false
	for _, branch := range branches {
		if branch == "develop" {
			developExist = true
		}
	}

	if developExist {
		git.Quiet("checkout", "develop")
	} else {
		git.Quiet("checkout", "-b", "develop")
	}

	if err != nil {
		return
	}

	_, err = git.Ref("HEAD")

	if err != nil {
		git.Run("commit", "-m", "Initial commit", "--allow-empty")
	}

	return
}

func FlowInit() (err error) {
	branching := getBranching()

	err = updateConfigFile(branching)

	if err != nil {
		return
	}

	err = initMasterBranch()
	err = initDevelopBranch()

	return
}

func CheckInit() {
	errorMessage := "You need to init your repository with git flow init before."
	gitConfig, err := NewConfig()

	if err != nil {
		utils.Check(fmt.Errorf("%s", err.Error()))
	}

	sectionBranch := gitConfig.GetSection(gitConfig.sectionBranch)

	if !sectionBranch.HasKey("develop") && !sectionBranch.HasKey("master") {
		utils.Check(fmt.Errorf("%s", errorMessage))
	}

	sectionPrefix := gitConfig.GetSection(gitConfig.sectionPrefix)

	if !sectionPrefix.HasKey("feature") ||
		!sectionPrefix.HasKey("release") ||
		!sectionPrefix.HasKey("hotfix") {
		utils.Check(fmt.Errorf("%s", errorMessage))
	}
}
