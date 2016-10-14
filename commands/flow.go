package commands

import (
	"fmt"

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
		/*err := flow.FlowFeatureStart(featureName)
		if err != nil {
			errorMessage = err.Error()
		}*/
		errorMessage = featureName
	case "finish":
		/*err := flow.FlowFeatureFinish(featureName)
		if err != nil {
			errorMessage = err.Error()
		}*/
		errorMessage = featureName
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
