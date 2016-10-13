package commands

import (
	"fmt"

	"github.com/github/hub/git"
	"github.com/github/hub/utils"
)

var (
	cmdFlow = &Command{
		Run:   flow,
		Usage: "flow ??? ???",
		Long: `Check out the head of a pull request as a local branch.

## Examples:
		$ hub checkout https://github.com/jingweno/gh/pull/73
		> git fetch origin pull/73/head:jingweno-feature
		> git checkout jingweno-feature

## See also:

hub-merge(1), hub-am(1), hub(1), git-checkout(1)
`,
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
	words := args.Words()
	if len(words) != 2 {
		utils.Check(fmt.Errorf("%s", cmdFlow.HelpText()))
	}

	instruction := words[0]
	featureName := words[1]

	switch instruction {
	case "start":
		flowFeatureStart(featureName)

	case "finish":
		fmt.Println("finish")
	default:
		fmt.Printf("%s", cmdFlow.HelpText())
	}
	args.NoForward()
}

func flowFeatureStart(featureName string) (err error) {
	branchName := "feature/" + featureName
	err = git.Spawn("checkout", "develop")

	if err == nil {
		git.Run("checkout", "-b", branchName)
	}

	return
}

func flow(command *Command, args *Args) {
	git.Run("status")
	args.NoForward()
}
