package commands

import (
	"fmt"

	gitFlow "github.com/boris-rea/hub/flow"
	"github.com/github/hub/utils"
)

var (
	cmdFlow = &Command{
		Run:   flow,
		Usage: `flow <type> <command> <name>`,
		Long: `Permit to use basic operations like in gitFlow

	## Examples:
		$ hub flow feature start myFeature
		[ Start a feature from develop and call it feature/myFeature ]

		$ hub flow feature finish myFeature
		[ Merge feature/myFeature in develop and delete the local feature branch ]

		$ hub flow release start myFeature
		[ Start a release from develop and call it release/myFeature ]

		$ hub flow release finish myFeature
		[ Merge release/myFeature in develop and master, tag the master branch and delete the local release branch ]

		$ hub flow hotfix start myFeature
		[ Start a hotfix from master and call it hotfix/myFeature ]

		$ hub flow hotfix finish myFeature
		[ Merge hotfix/myFeature in develop and master, tag the master branch and delete the local release branch ]
`,
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

var (
	flagCreatePullRequest bool
)

func init() {
	cmdFlowFeature.Flag.BoolVarP(&flagCreatePullRequest, "pull-request", "", false, "PULLREQUEST")

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
		var err error
		if flagCreatePullRequest {
			err = gitFlow.FlowFeaturePullRequest(featureName)
		} else {
			err = gitFlow.FlowFeatureFinish(featureName)
		}

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
