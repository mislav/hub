package commands

import (
	"fmt"

	"github.com/github/hub/git"
	"github.com/github/hub/utils"
)

var (
	cmdFlow = &Command{
		Run:   flow,
		Usage: `flow`,
		Long:  `TODO`,
	}

	cmdFlowFeature = &Command{
		Key: "feature",
		Run: flowFeature,
	}
)

func init() {
	cmdFlow.Use(cmdFlowFeature)
	CmdRunner.Use(cmdFlow)
}

func flowFeature(command *Command, args *Args) {
	args.NoForward()
	words := args.Words()

	if len(words) != 2 {
		utils.Check(fmt.Errorf("%s", cmdFlow.HelpText()))
	}

	errorMessage := ""
	instruction := words[0]
	featureName := words[1]

	switch instruction {
	case "start":
		err := flowFeatureStart(featureName)
		if err != nil {
			errorMessage = err.Error()
		}
	case "finish":
		err := flowFeatureFinish(featureName)
		if err != nil {
			errorMessage = err.Error()
		}
	default:
		errorMessage = cmdFlow.HelpText()
	}

	if errorMessage != "" {
		utils.Check(fmt.Errorf("%s", errorMessage))
	}
}

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
	err = git.Spawn("checkout", "develop")

	if err == nil {
		err = git.Spawn("merge", branchName)
	}

	if err == nil {
		err = git.Spawn("branch", "-d", branchName)
	}

	return
}

func flow(command *Command, args *Args) {
	args.NoForward()
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
