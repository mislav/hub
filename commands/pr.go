package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var (
	cmdPr = &Command{
		Run: printHelp,
		Usage: `
pr list [-s <STATE>] [-h <HEAD>] [-b <BASE>] [-o <SORT_KEY> [-^]] [-f <FORMAT>] [-L <LIMIT>]
pr checkout <PR-NUMBER> [<BRANCH>]
`,
		Long: `Manage GitHub pull requests for the current project.

## Commands:

	* _list_:
		List pull requests in the current project.

	* _checkout_:
		Check out the head of a pull request in a new branch.

## Options:

	-s, --state <STATE>
		Filter pull requests by <STATE>. Supported values are: "open" (default),
		"closed", "merged", or "all".

	-h, --head <BRANCH>
		Show pull requests started from the specified head <BRANCH>. The
		"OWNER:BRANCH" format must be used for pull requests from forks.

	-b, --base <BRANCH>
		Show pull requests based off the specified <BRANCH>.

	-f, --format <FORMAT>
		Pretty print the list of pull requests using format <FORMAT> (default:
		"%sC%>(8)%i%Creset  %t%  l%n"). See the "PRETTY FORMATS" section of
		git-log(1) for some additional details on how placeholders are used in
		format. The available placeholders are:

		%I: pull request number

		%i: pull request number prefixed with "#"

		%U: the URL of this pull request

		%S: state ("open" or "closed")

		%sC: set color to red or green, depending on pull request state.

		%t: title

		%l: colored labels

		%L: raw, comma-separated labels

		%b: body

		%B: base branch

		%sB: base commit SHA

		%H: head branch

		%sH: head commit SHA

		%sm: merge commit SHA

		%au: login name of author

		%as: comma-separated list of assignees

		%rs: comma-separated list of requested reviewers

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

		%mD: merged date-only (no time of day)

		%mr: merged date, relative

		%mt: merged date, UNIX timestamp

		%mI: merged date, ISO 8601 format

		%n: newline

		%%: a literal %

	-o, --sort <KEY>
		Sort displayed issues by "created" (default), "updated", "popularity", or "long-running".

	-^, --sort-ascending
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
		Key:  "list",
		Run:  listPulls,
		Long: cmdPr.Long,
	}
)

func init() {
	cmdPr.Use(cmdListPulls)
	cmdPr.Use(cmdCheckoutPr)
	CmdRunner.Use(cmdPr)
}

func printHelp(command *Command, args *Args) {
	utils.Check(command.UsageError(""))
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

	filters := map[string]interface{}{}
	if args.Flag.HasReceived("--state") {
		filters["state"] = args.Flag.Value("--state")
	}
	if args.Flag.HasReceived("--sort") {
		filters["sort"] = args.Flag.Value("--sort")
	}
	if args.Flag.HasReceived("--base") {
		filters["base"] = args.Flag.Value("--base")
	}
	if args.Flag.HasReceived("--head") {
		head := args.Flag.Value("--head")
		if !strings.Contains(head, ":") {
			head = fmt.Sprintf("%s:%s", project.Owner, head)
		}
		filters["head"] = head
	}

	if args.Flag.Bool("--sort-ascending") {
		filters["direction"] = "asc"
	} else {
		filters["direction"] = "desc"
	}

	onlyMerged := false
	if filters["state"] == "merged" {
		filters["state"] = "closed"
		onlyMerged = true
	}

	flagPullRequestLimit := args.Flag.Int("--limit")
	flagPullRequestFormat := args.Flag.Value("--format")
	if !args.Flag.HasReceived("--format") {
		flagPullRequestFormat = "%sC%>(8)%i%Creset  %t%  l%n"
	}

	pulls, err := gh.FetchPullRequests(project, filters, flagPullRequestLimit, func(pr *github.PullRequest) bool {
		return !(onlyMerged && pr.MergedAt.IsZero())
	})
	utils.Check(err)

	colorize := ui.IsTerminal(os.Stdout)
	for _, pr := range pulls {
		ui.Print(formatPullRequest(pr, flagPullRequestFormat, colorize))
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
	placeholders := formatIssuePlaceholders(github.Issue(pr), colorize)
	for key, value := range formatPullRequestPlaceholders(pr) {
		placeholders[key] = value
	}
	return ui.Expand(format, placeholders, colorize)
}
