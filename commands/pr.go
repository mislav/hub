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

/**
 * Add a log messsage to /tmp/hub.log
 *
 * FIXME: Should be removed before PRing this.
 */
func johan(message string) {
	f, err := os.OpenFile("/tmp/hub.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(message + "\n"); err != nil {
		panic(err)
	}
}

func checkoutPr(command *Command, args *Args) {
	johan("")

	if args.ParamsSize() < 1 || args.ParamsSize() > 2 {
		utils.Check(fmt.Errorf("Error: Expected one or two arguments, got %d", args.ParamsSize()))
	}

	prNumberString := args.GetParam(0)
	_, err := strconv.Atoi(prNumberString)
	utils.Check(err)

	johan("ormen")

	// Figure out the PR URL
	localRepo, err := github.LocalRepo()
	johan("eva")
	utils.Check(err)
	johan("adam")
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
