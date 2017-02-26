package commands

import (
	"fmt"
	"strconv"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var cmdPr = &Command{
	Run:   pr,
	Usage: "pr <PULLREQ-NUMBER> [<BRANCH>]",
	Long: `Check out the head of a pull request as a local branch.

## Examples:
		$ hub pr 73
		> git fetch origin pull/73/head:jingweno-feature
		> git checkout jingweno-feature

## See also:

hub-merge(1), hub(1), hub-checkout(1)
`,
}

func init() {
	CmdRunner.Use(cmdPr)
}

func pr(command *Command, args *Args) {
	if args.ParamsSize() < 1 || args.ParamsSize() > 2 {
		utils.Check(fmt.Errorf("Error: Expected one or two arguments, got %d", args.ParamsSize()))
	}

	prNumberString := args.GetParam(0)
	_, err := strconv.Atoi(prNumberString)
	utils.Check(err)

	// Figure out the PR URL
	localRepo, err := github.LocalRepo()
	utils.Check(err)
	baseProject, err := localRepo.MainProject()
	utils.Check(err)
	host, err := github.CurrentConfig().PromptForHost(baseProject.Host)
	utils.Check(err)
	client := github.NewClientWithHost(host)
	pr, err := client.PullRequest(baseProject, prNumberString)
	utils.Check(err)

	if args.ParamsSize() == 1 {
		args.Replace(args.Executable, "checkout", pr.HtmlUrl)
	} else {
		args.Replace(args.Executable, "checkout", pr.HtmlUrl, args.GetParam(1))
	}

	// Call into the checkout code which already provides the functionality we're
	// after
	err = transformCheckoutArgs(args)
	utils.Check(err)
}
