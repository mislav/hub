package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var (
	cmdPr = &Command{
		Run: printHelp,
		Usage: `
pr list [-s <STATE>] [-h <HEAD>] [-b <BASE>] [-o <SORT_KEY> [-^]] [-L <LIMIT>]
pr checkout <PR-NUMBER> [<BRANCH>]
pr push <PULLREQ-NUMBER> [<BRANCH>]
`,
		Long: `Manage GitHub pull requests for the current project.

## Commands:

	* _list_:
		List pull requests in the current project.

	* _checkout_:
		Check out the head of a pull request in a new branch.

	* _push_:
		Push a local branch to the head of a pull request.

## Options:

	-s, --state <STATE>
		Display pull requests with state <STATE> (default: "open").

	-f, --format <FORMAT>
		Pretty print the list of pull requests using format <FORMAT> (default:
		"%sC%>(8)%i%Creset  %t%  l%n"). See the "PRETTY FORMATS" section of the
		git-log manual for some additional details on how placeholders are used in
		format. The available placeholders are:

		%I: pull request number

		%i: pull request number prefixed with "#"

		%U: the URL of this pull request

		%S: state (i.e. "open", "closed")

		%sC: set color to red or green, depending on pull request state.

		%t: title

		%l: colored labels

		%L: raw, comma-separated labels

		%b: body

		%B: base branch

		%H: head branch

		%au: login name of author

		%as: comma-separated list of assignees

		%Mn: milestone number

		%Mt: milestone title

		%NC: number of comments

		%Nc: number of comments wrapped in parentheses, or blank string if zero.

		%cD: created date-only (no time of day)

		%cr: created date, relative

		%ct: created date, UNIX timestamp

		%cI: created date, ISO 8601 format

		%uD: updated date-only (no time of day)

		%ur: updated date, relative

		%ut: updated date, UNIX timestamp

		%uI: updated date, ISO 8601 format

	-o, --sort <SORT_KEY>
		Sort displayed issues by "created" (default), "updated", "popularity", or "long-running".

	-^ --sort-ascending
		Sort by ascending dates instead of descending.

	-L, --limit <LIMIT>
		Display only the first <LIMIT> issues.

## See also:

hub-issue(1), hub-pull-request(1), hub(1)
`,
	}

	cmdCheckoutPr = &Command{
		Key: "checkout",
		Run: checkoutPr,
	}

	cmdListPulls = &Command{
		Key: "list",
		Run: listPulls,
	}

	cmdPushPr = &Command{
		Key: "push",
		Run: pushPr,
	}

	flagPullRequestState,
	flagPullRequestFormat,
	flagPullRequestSort string

	flagPullRequestSortAscending bool

	flagPullRequestLimit int

	flagForce,
	flagSetUpstream bool
)

func init() {
	cmdListPulls.Flag.StringVarP(&flagPullRequestState, "state", "s", "", "STATE")
	cmdListPulls.Flag.StringVarP(&flagPullRequestBase, "base", "b", "", "BASE")
	cmdListPulls.Flag.StringVarP(&flagPullRequestHead, "head", "h", "", "HEAD")
	cmdListPulls.Flag.StringVarP(&flagPullRequestFormat, "format", "f", "%sC%>(8)%i%Creset  %t%  l%n", "FORMAT")
	cmdListPulls.Flag.StringVarP(&flagPullRequestSort, "sort", "o", "created", "SORT_KEY")
	cmdListPulls.Flag.BoolVarP(&flagPullRequestSortAscending, "sort-ascending", "^", false, "SORT_KEY")
	cmdListPulls.Flag.IntVarP(&flagPullRequestLimit, "limit", "L", -1, "LIMIT")

	cmdPr.Use(cmdListPulls)
	cmdPr.Use(cmdCheckoutPr)

	cmdPushPr.Flag.BoolVarP(&flagForce, "force", "f", false, "FORCE")
	cmdPushPr.Flag.BoolVarP(&flagSetUpstream, "set-upstream", "u", false, "SET-UPSTREAM")
	cmdPr.Use(cmdPushPr)

	CmdRunner.Use(cmdPr)
}

func printHelp(command *Command, args *Args) {
	fmt.Print(command.HelpText())
	os.Exit(0)
}

func listPulls(cmd *Command, args *Args) {
	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.MainProject()
	utils.Check(err)

	gh := github.NewClient(project.Host)

	args.NoForward()
	if args.Noop {
		ui.Printf("Would request list of pull requests for %s\n", project)
		return
	}

	flagFilters := map[string]string{
		"state": flagPullRequestState,
		"head":  flagPullRequestHead,
		"base":  flagPullRequestBase,
		"sort":  flagPullRequestSort,
	}
	filters := map[string]interface{}{}
	for flag, filter := range flagFilters {
		if cmd.FlagPassed(flag) {
			filters[flag] = filter
		}
	}
	if flagPullRequestSortAscending {
		filters["direction"] = "asc"
	}

	pulls, err := gh.FetchPullRequests(project, filters, flagPullRequestLimit, nil)
	utils.Check(err)

	colorize := ui.IsTerminal(os.Stdout)
	for _, pr := range pulls {
		ui.Printf(formatPullRequest(pr, flagPullRequestFormat, colorize))
	}
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

func formatPullRequest(pr github.PullRequest, format string, colorize bool) string {
	base := pr.Base.Ref
	head := pr.Head.Label
	if pr.IsSameRepo() {
		head = pr.Head.Ref
	}

	placeholders := formatIssuePlaceholders(github.Issue(pr), colorize)
	placeholders["B"] = base
	placeholders["H"] = head

	return ui.Expand(format, placeholders, colorize)
}

func pushPr(command *Command, args *Args) {
	var localBranchName string
	words := args.Words()

	localRepo, err := github.LocalRepo()
	utils.Check(err)

	if len(words) == 0 {
		utils.Check(fmt.Errorf("Error: No pull request number given"))
	} else if len(words) > 1 {
		localBranchName = words[1]
	} else {
		branch, err := localRepo.CurrentBranch()
		utils.Check(err)
		localBranchName = branch.ShortName()
	}

	prNumberString := words[0]
	_, err = strconv.Atoi(prNumberString)
	utils.Check(err)

	baseProject, err := localRepo.MainProject()
	utils.Check(err)
	host, err := github.CurrentConfig().PromptForHost(baseProject.Host)
	utils.Check(err)
	client := github.NewClientWithHost(host)

	pr, err := client.PullRequest(baseProject, prNumberString)
	utils.Check(err)
	head := pr.Head
	headRepo := head.Repo

	refspec := fmt.Sprintf("%s:%s", localBranchName, head.Ref)

	remoteProject := github.NewProject(headRepo.Owner.Login, headRepo.Name, host.Host)
	remote := remoteProject.GitURL("", "", true)

	args.Replace(args.Executable, "push", remote, refspec, "--no-follow-tags")
	if flagForce {
		args.AppendParams("--force")
	}
	if flagSetUpstream {
		args.AppendParams("--set-upstream")
	}
}
