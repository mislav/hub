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
		Run: printHelp,
		Usage: `
pr checkout <PULLREQ-NUMBER> [<BRANCH>]
pr show [-b <BASE>] [-h <HEAD>]
pr help <COMMAND>
`,
	}

	cmdCheckoutPr = &Command{
		Key:   "checkout",
		Run:   checkoutPr,
		Usage: `pr checkout <PULLREQ-NUMBER> [<BRANCH>]`,
		Long: `Check out the head of a pull request as a local branch.

## Examples:
	$ hub pr checkout 73
	> git fetch origin pull/73/head:jingweno-feature
	> git checkout jingweno-feature

## See also:

hub-merge(1), hub(1), hub-checkout(1)
	`,
	}

	cmdShowPr = &Command{
		Key:   "show",
		Run:   showPr,
		Usage: `show [-b <BASE>] [-h <HEAD>]`,
		Long: `Display information about the pull request for the current branch.

## Examples:
	$ hub pr show
	$ hub pr show -b devel

## See also:

hub-pull-request(1), hub(1)
	`,
	}

	cmdHelpPr = &Command{
		Key: "help",
		Run: helpPr,
	}

	flagPullRequestShowBase,
	flagPullRequestShowHead string
)

func init() {
	cmdShowPr.Flag.StringVarP(&flagPullRequestShowBase, "base", "b", "", "BASE")
	cmdShowPr.Flag.StringVarP(&flagPullRequestShowHead, "head", "h", "", "HEAD")

	cmdPr.Use(cmdCheckoutPr)
	cmdPr.Use(cmdShowPr)
	cmdPr.Use(cmdHelpPr)
	CmdRunner.Use(cmdPr)
}

func printHelp(command *Command, args *Args) {
	fmt.Print(command.HelpText())
	os.Exit(0)
}

func helpPr(command *Command, args *Args) {
	words := args.Words()
	if len(words) == 0 {
		fmt.Print(cmdPr.HelpText())
	} else {
		switch words[0] {
		case "checkout":
			fmt.Print(cmdCheckoutPr.HelpText())
		case "show":
			fmt.Print(cmdShowPr.HelpText())
		default:
			utils.Check(fmt.Errorf("Error: No such subcommand: %s", words[0]))
		}
	}
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

func showPr(command *Command, args *Args) {
	localRepo, err := github.LocalRepo()
	utils.Check(err)

	currentBranch, err := localRepo.CurrentBranch()
	utils.Check(err)

	baseProject, err := localRepo.MainProject()
	utils.Check(err)

	host, err := github.CurrentConfig().PromptForHost(baseProject.Host)
	utils.Check(err)
	client := github.NewClientWithHost(host)

	trackedBranch, headProject, err := localRepo.RemoteBranchAndProject(host.User, false)
	utils.Check(err)

	var base, head string

	if flagPullRequestShowBase != "" {
		baseProject, base = parsePullRequestProject(baseProject, flagPullRequestShowBase)
	}

	if flagPullRequestHead != "" {
		headProject, head = parsePullRequestProject(headProject, flagPullRequestShowHead)
	}

	if base == "" {
		masterBranch := localRepo.MasterBranch()
		base = masterBranch.ShortName()
	}

	if head == "" && trackedBranch != nil {
		if !trackedBranch.IsRemote() {
			// the current branch tracking another branch
			// pretend there's no upstream at all
			trackedBranch = nil
		} else {
			if baseProject.SameAs(headProject) && base == trackedBranch.ShortName() {
				e := fmt.Errorf(`Aborted: head branch is the same as base ("%s")`, base)
				e = fmt.Errorf("%s\n(use `-h <branch>` to specify an explicit pull request head)", e)
				utils.Check(e)
			}
		}
	}

	if head == "" {
		if trackedBranch != nil {
			head = currentBranch.ShortName()
		} else {
			head = trackedBranch.ShortName()
		}
	}

	if headRepo, err := client.Repository(headProject); err == nil {
		headProject.Owner = headRepo.Owner.Login
		headProject.Name = headRepo.Name
	}

	filters := map[string]interface{}{
		"base":  base,
		"head":  fmt.Sprintf("%s:%s", headProject.Owner, head),
		"state": "open",
	}
	prs, err := client.FetchPullRequests(baseProject, filters)
	utils.Check(err)

	if len(prs) == 0 {
		utils.Check(fmt.Errorf("No pull requests found"))
		os.Exit(1)
	} else {
		for _, pr := range prs {
			fmt.Println(pr.HtmlUrl)
		}
		os.Exit(0)
	}
}
