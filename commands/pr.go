package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/github/hub/v2/git"
	"github.com/github/hub/v2/github"
	"github.com/github/hub/v2/ui"
	"github.com/github/hub/v2/utils"
)

var (
	cmdPr = &Command{
		Run: printHelp,
		Usage: `
pr list [-s <STATE>] [-h <HEAD>] [-b <BASE>] [-o <SORT_KEY> [-^]] [-f <FORMAT>] [-L <LIMIT>]
pr checkout <PR-NUMBER> [<BRANCH>]
pr show [-uc] [-f <FORMAT>] [-h <HEAD>]
pr show [-uc] [-f <FORMAT>] <PR-NUMBER>
pr merge [-d] [--squash | --rebase] <PR-NUMBER> [-m <MESSAGE> | -F <FILE>] [--head-sha <COMMIT-SHA>]
`,
		Long: `Manage GitHub Pull Requests for the current repository.

## Commands:

	* _list_:
		List pull requests in the current repository.

	* _checkout_:
		Check out the head of a pull request in a new branch.

		To update the pull request with new commits, use ''git push''.

	* _show_:
		Open a pull request page in a web browser. When no <PR-NUMBER> is
		specified, <HEAD> is used to look up open pull requests and defaults to
		the current branch name. With ''--format'', print information about the
		pull request instead of opening it.

	* _merge_:
		Merge a pull request in the current repository remotely. Select an
		alternate merge method with ''--squash'' or ''--rebase''. Change the
		commit subject and body with ''--message'' or ''--file''.

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
		"%pC%>(8)%i%Creset  %t%  l%n"). See the "PRETTY FORMATS" section of
		git-log(1) for some additional details on how placeholders are used in
		format. The available placeholders are:

		%I: pull request number

		%i: pull request number prefixed with "#"

		%U: the URL of this pull request

		%S: state ("open" or "closed")

		%pS: pull request state ("open", "draft", "merged", or "closed")

		%sC: set color to red or green, depending on state

		%pC: set color according to pull request state

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

	--color[=<WHEN>]
		Enable colored output even if stdout is not a terminal. <WHEN> can be one
		of "always" (default for ''--color''), "never", or "auto" (default).

	-o, --sort <KEY>
		Sort displayed pull requests by "created" (default), "updated", "popularity", or "long-running".

	-^, --sort-ascending
		Sort by ascending dates instead of descending.

	-L, --limit <LIMIT>
		Display only the first <LIMIT> pull requests.

	-u, --url
		Print the pull request URL instead of opening it.

	-c, --copy
		Put the pull request URL to clipboard instead of opening it.

	-m, --message <MESSAGE>
		The text up to the first blank line in <MESSAGE> is treated as the commit
		subject for the merge commit, and the rest is used as commit body.

		When multiple ''--message'' are passed, their values are concatenated with a
		blank line in-between.

	-F, --file <FILE>
		Read the subject and body for the merge commit from <FILE>. Pass "-" to read
		from standard input instead. See ''--message'' for the formatting rules.

	--head-sha <COMMIT-SHA>
		Ensure that the head of the pull request matches the commit SHA when merging.

	--squash
		Squash commits instead of creating a merge commit when merging a pull request.

	--rebase
		Rebase commits on top of the base branch when merging a pull request.

	-d, --delete-branch
		Delete the head branch after successfully merging a pull request.

## See also:

hub-issue(1), hub-pull-request(1), hub(1)
`,
	}

	cmdCheckoutPr = &Command{
		Key:        "checkout",
		Run:        checkoutPr,
		KnownFlags: "\n",
	}

	cmdListPulls = &Command{
		Key:  "list",
		Run:  listPulls,
		Long: cmdPr.Long,
	}

	cmdShowPr = &Command{
		Key: "show",
		Run: showPr,
		KnownFlags: `
		-h, --head HEAD
		-u, --url
		-c, --copy
		-f, --format FORMAT
		--color
		`,
	}

	cmdMergePr = &Command{
		Key: "merge",
		Run: mergePr,
		KnownFlags: `
		-m, --message MESSAGE
		-F, --file FILE
		--head-sha COMMIT
		--squash
		--rebase
		-d, --delete-branch
		`,
	}
)

func init() {
	cmdPr.Use(cmdListPulls)
	cmdPr.Use(cmdCheckoutPr)
	cmdPr.Use(cmdShowPr)
	cmdPr.Use(cmdMergePr)
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
		flagPullRequestFormat = "%pC%>(8)%i%Creset  %t%  l%n"
	}

	pulls, err := gh.FetchPullRequests(project, filters, flagPullRequestLimit, func(pr *github.PullRequest) bool {
		return !(onlyMerged && pr.MergedAt.IsZero())
	})
	utils.Check(err)

	colorize := colorizeOutput(args.Flag.HasReceived("--color"), args.Flag.Value("--color"))
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

func showPr(command *Command, args *Args) {
	localRepo, err := github.LocalRepo()
	utils.Check(err)

	baseProject, err := localRepo.MainProject()
	utils.Check(err)

	host, err := github.CurrentConfig().PromptForHost(baseProject.Host)
	utils.Check(err)
	gh := github.NewClientWithHost(host)

	words := args.Words()
	openURL := ""
	prNumber := 0
	var pr *github.PullRequest

	if len(words) > 0 {
		if prNumber, err = strconv.Atoi(words[0]); err == nil {
			openURL = baseProject.WebURL("", "", fmt.Sprintf("pull/%d", prNumber))
		} else {
			utils.Check(fmt.Errorf("invalid pull request number: '%s'", words[0]))
		}
	} else {
		pr, err = findCurrentPullRequest(localRepo, gh, baseProject, args.Flag.Value("--head"))
		utils.Check(err)
		openURL = pr.HTMLURL
	}

	args.NoForward()
	if format := args.Flag.Value("--format"); format != "" {
		if pr == nil {
			pr, err = gh.PullRequest(baseProject, strconv.Itoa(prNumber))
			utils.Check(err)
		}
		colorize := colorizeOutput(args.Flag.HasReceived("--color"), args.Flag.Value("--color"))
		ui.Println(formatPullRequest(*pr, format, colorize))
		return
	}

	printURL := args.Flag.Bool("--url")
	copyURL := args.Flag.Bool("--copy")

	printBrowseOrCopy(args, openURL, !printURL && !copyURL, copyURL)
}

func findCurrentPullRequest(localRepo *github.GitHubRepo, gh *github.Client, baseProject *github.Project, headArg string) (*github.PullRequest, error) {
	filterParams := map[string]interface{}{
		"state": "open",
	}
	headWithOwner := ""

	if headArg != "" {
		headWithOwner = headArg
		if !strings.Contains(headWithOwner, ":") {
			headWithOwner = fmt.Sprintf("%s:%s", baseProject.Owner, headWithOwner)
		}
	} else {
		currentBranch, err := localRepo.CurrentBranch()
		utils.Check(err)
		if headBranch, headProject, err := findPushTarget(currentBranch); err == nil {
			headWithOwner = fmt.Sprintf("%s:%s", headProject.Owner, headBranch.ShortName())
		} else if headProject, err := deducePushTarget(currentBranch, gh.Host.User); err == nil {
			headWithOwner = fmt.Sprintf("%s:%s", headProject.Owner, currentBranch.ShortName())
		} else {
			headWithOwner = fmt.Sprintf("%s:%s", baseProject.Owner, currentBranch.ShortName())
		}
	}

	filterParams["head"] = headWithOwner

	pulls, err := gh.FetchPullRequests(baseProject, filterParams, 1, nil)
	if err != nil {
		return nil, err
	} else if len(pulls) == 1 {
		return &pulls[0], nil
	} else {
		return nil, fmt.Errorf("no open pull requests found for branch '%s'", headWithOwner)
	}
}

func branchTrackingInformation(branch *github.Branch) (string, *github.Branch, error) {
	branchRemote, err := git.Config(fmt.Sprintf("branch.%s.remote", branch.ShortName()))
	if branchRemote == "." {
		err = fmt.Errorf("branch is tracking another local branch")
	}
	if err != nil {
		return "", nil, err
	}
	branchMerge, err := git.Config(fmt.Sprintf("branch.%s.merge", branch.ShortName()))
	if err != nil {
		return "", nil, err
	}
	trackingBranch := &github.Branch{
		Repo: branch.Repo,
		Name: branchMerge,
	}
	return branchRemote, trackingBranch, nil
}

func findPushTarget(branch *github.Branch) (*github.Branch, *github.Project, error) {
	branchRemote, headBranch, err := branchTrackingInformation(branch)
	if err != nil {
		return nil, nil, err
	}

	if headRemote, err := branch.Repo.RemoteByName(branchRemote); err == nil {
		headProject, err := headRemote.Project()
		if err != nil {
			return nil, nil, err
		}
		return headBranch, headProject, nil
	}

	remoteURL, err := git.ParseURL(branchRemote)
	if err != nil {
		return nil, nil, err
	}
	headProject, err := github.NewProjectFromURL(remoteURL)
	if err != nil {
		return nil, nil, err
	}
	return headBranch, headProject, nil
}

func deducePushTarget(branch *github.Branch, owner string) (*github.Project, error) {
	remote := branch.Repo.RemoteForBranch(branch, owner)
	if remote == nil {
		return nil, fmt.Errorf("no remote found for branch %s", branch.ShortName())
	}
	return remote.Project()
}

func mergePr(command *Command, args *Args) {
	words := args.Words()
	if len(words) == 0 {
		utils.Check(fmt.Errorf("Error: No pull request number given"))
	}

	prNumber, err := strconv.Atoi(words[0])
	utils.Check(err)

	params := map[string]interface{}{
		"merge_method": "merge",
	}
	if args.Flag.Bool("--squash") {
		params["merge_method"] = "squash"
	}
	if args.Flag.Bool("--rebase") {
		params["merge_method"] = "rebase"
	}

	msgs := args.Flag.AllValues("--message")
	if len(msgs) > 0 {
		params["commit_title"] = msgs[0]
		params["commit_message"] = strings.Join(msgs[1:], "\n\n")
	} else if args.Flag.HasReceived("--file") {
		content, err := msgFromFile(args.Flag.Value("--file"))
		utils.Check(err)
		params["commit_title"], params["commit_message"] = github.SplitTitleBody(content)
	}

	if headSHA := args.Flag.Value("--head-sha"); headSHA != "" {
		params["sha"] = args.Flag.Value("--head-sha")
	}

	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.MainProject()
	utils.Check(err)

	args.NoForward()
	if args.Noop {
		ui.Printf("Would merge pull request #%d for %s\n", prNumber, project)
		return
	}

	gh := github.NewClient(project.Host)
	_, err = gh.MergePullRequest(project, prNumber, params)
	utils.Check(err)

	if !args.Flag.Bool("--delete-branch") {
		return
	}

	pr, err := gh.PullRequest(project, strconv.Itoa(prNumber))
	utils.Check(err)
	if !pr.IsSameRepo() {
		return
	}

	branchName := pr.Head.Ref
	err = gh.DeleteBranch(project, branchName)
	utils.Check(err)
}

func formatPullRequest(pr github.PullRequest, format string, colorize bool) string {
	placeholders := formatIssuePlaceholders(github.Issue(pr), colorize)
	delete(placeholders, "NC")
	delete(placeholders, "Nc")

	for key, value := range formatPullRequestPlaceholders(pr, colorize) {
		placeholders[key] = value
	}
	return ui.Expand(format, placeholders, colorize)
}
