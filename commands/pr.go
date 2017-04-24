package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var (
	cmdPr = &Command{
		Run:   printHelp,
		Usage: "pr checkout <PULLREQ-NUMBER> [<BRANCH>]",
		Long: `Check out the head of a pull request as a local branch.

## Examples:
	$ hub pr checkout 73
	> git fetch origin pull/73/head:jingweno-feature
	> git checkout jingweno-feature

## See also:

hub-merge(1), hub(1), hub-checkout(1)
	`,
	}

	cmdCheckoutPr = &Command{
		Key: "checkout",
		Run: checkoutPr,
	}
)

func init() {
	cmdPr.Use(cmdCheckoutPr)
	CmdRunner.Use(cmdPr)
}

func printHelp(command *Command, args *Args) {
	fmt.Print(command.HelpText())
	os.Exit(0)
}

func checkoutPr(command *Command, args *Args) {
	words := args.Words()
	var newBranchName string

	if len(words) == 0 {
		utils.Check(fmt.Errorf("Error: No pull request number given"))
	} else if len(words) > 1 {
		newBranchName = words[1]
	}

	prNumberString := words[0]
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

	newArgs, err := transformCheckoutArgs(args, pr, newBranchName)
	utils.Check(err)

	args.Replace(args.Executable, "checkout", newArgs...)
}
