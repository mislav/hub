package commands

import (
	"fmt"

	gitFlow "github.com/boris-rea/hub/flow"
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

	cmdFlowRelease = &Command{
		Key: "release",
		Run: flowRelease,
	}

	cmdFlowHotfix = &Command{
		Key: "hotfix",
		Run: flowHotfix,
	}
)

func init() {
	cmdFlow.Use(cmdFlowFeature)
	cmdFlow.Use(cmdFlowRelease)
	cmdFlow.Use(cmdFlowHotfix)
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
		err := gitFlow.FlowFeatureStart(featureName)
		if err != nil {
			errorMessage = err.Error()
		}
	case "finish":
		err := gitFlow.FlowFeatureFinish(featureName)
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

func flowRelease(command *Command, args *Args) {
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
		err := gitFlow.FlowReleaseStart(featureName)
		if err != nil {
			errorMessage = err.Error()
		}
	case "finish":
		err := gitFlow.FlowReleaseFinish(featureName)
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

func flowHotfix(command *Command, args *Args) {
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
		err := gitFlow.FlowHotfixStart(featureName)
		if err != nil {
			errorMessage = err.Error()
		}
	case "finish":
		err := gitFlow.FlowHotfixFinish(featureName)
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

func flow(command *Command, args *Args) {
	args.NoForward()
}
