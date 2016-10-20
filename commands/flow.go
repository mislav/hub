package commands

import (
	"fmt"

	gitFlow "github.com/boris-rea/hub/flow"
	"github.com/github/hub/utils"
)

var (
	cmdFlow = &Command{
		Run:   flow,
		Usage: `flow <type> [<command> <name>]`,
		Long: `Permit to use basic operations like in gitFlow

## Examples:
	$ hub flow init
	[ Init config variable for branching ]

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

	cmdFlowInit = &Command{
		Key: "init",
		Run: flowInit,
	}
)

var (
	flagPullRequestAssigneesFlow listFlag
	flagCreatePullRequest        bool
)

func init() {
	cmdFlowFeature.Flag.BoolVarP(&flagCreatePullRequest, "pull-request", "", false, "PULLREQUESTFEATURE")
	cmdFlowFeature.Flag.VarP(&flagPullRequestAssignees, "assign", "a", "")
	cmdFlowRelease.Flag.BoolVarP(&flagCreatePullRequest, "pull-request", "", false, "PULLREQUESTRELEASE")
	cmdFlowRelease.Flag.VarP(&flagPullRequestAssignees, "assign", "a", "USERS")
	cmdFlowHotfix.Flag.BoolVarP(&flagCreatePullRequest, "pull-request", "", false, "PULLREQUESTHOTFIX")
	cmdFlowHotfix.Flag.VarP(&flagPullRequestAssignees, "assign", "a", "USERS")

	cmdFlow.Use(cmdFlowFeature)
	cmdFlow.Use(cmdFlowRelease)
	cmdFlow.Use(cmdFlowHotfix)
	cmdFlow.Use(cmdFlowInit)
	CmdRunner.Use(cmdFlow)
}

func flowFeature(command *Command, args *Args) {
	args.NoForward()
	gitFlow.CheckInit()
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
			params := map[string]string{}
			if len(flagPullRequestAssignees) > 0 {
				params["assignees"] = flagPullRequestAssignees.String()
			}
			err = gitFlow.FlowFeaturePullRequest(featureName, params)
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
	gitFlow.CheckInit()
	words := args.Words()

	if len(words) != 2 {
		utils.Check(fmt.Errorf("%s", cmdFlow.HelpText()))
	}

	errorMessage := ""
	instruction := words[0]
	releaseName := words[1]

	switch instruction {
	case "start":
		err := gitFlow.FlowReleaseStart(releaseName)
		if err != nil {
			errorMessage = err.Error()
		}
	case "finish":
		var err error
		if flagCreatePullRequest {
			params := map[string]string{}
			if len(flagPullRequestAssignees) > 0 {
				params["assignees"] = flagPullRequestAssignees.String()
			}
			err = gitFlow.FlowReleasePullRequest(releaseName, params)
		} else {
			err = gitFlow.FlowReleaseFinish(releaseName)
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

func flowHotfix(command *Command, args *Args) {
	args.NoForward()
	gitFlow.CheckInit()
	words := args.Words()

	if len(words) != 2 {
		utils.Check(fmt.Errorf("%s", cmdFlow.HelpText()))
	}

	errorMessage := ""
	instruction := words[0]
	hotfixName := words[1]

	switch instruction {
	case "start":
		err := gitFlow.FlowHotfixStart(hotfixName)
		if err != nil {
			errorMessage = err.Error()
		}
	case "finish":
		var err error
		if flagCreatePullRequest {
			params := map[string]string{}
			if len(flagPullRequestAssignees) > 0 {
				params["assignees"] = flagPullRequestAssignees.String()
			}
			err = gitFlow.FlowHotfixPullRequest(hotfixName, params)
		} else {
			err = gitFlow.FlowHotfixFinish(hotfixName)
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

func flowInit(command *Command, args *Args) {
	args.NoForward()
	words := args.Words()

	if len(words) > 0 {
		utils.Check(fmt.Errorf("%s", cmdFlow.HelpText()))
	}

	err := gitFlow.FlowInit()

	if err != nil {
		utils.Check(fmt.Errorf("%s", err.Error()))
	}
}

func flow(command *Command, args *Args) {
	args.NoForward()
	gitFlow.CheckInit()
	utils.Check(fmt.Errorf("%s", command.HelpText()))
}
