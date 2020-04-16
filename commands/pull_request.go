package commands

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/github/hub/v2/git"
	"github.com/github/hub/v2/github"
	"github.com/github/hub/v2/utils"
)

var cmdPullRequest = &Command{
	Run: pullRequest,
	Usage: `
pull-request [-focpd] [-b <BASE>] [-h <HEAD>] [-r <REVIEWERS> ] [-a <ASSIGNEES>] [-M <MILESTONE>] [-l <LABELS>]
pull-request -m <MESSAGE> [--edit]
pull-request -F <FILE> [--edit]
pull-request -i <ISSUE>
`,
	Long: `Create a GitHub Pull Request.

## Options:
	-f, --force
		Skip the check for unpushed commits.

	-m, --message <MESSAGE>
		The text up to the first blank line in <MESSAGE> is treated as the pull
		request title, and the rest is used as pull request description in Markdown
		format.

		When multiple ''--message'' are passed, their values are concatenated with a
		blank line in-between.

		When neither ''--message'' nor ''--file'' were supplied, a text editor will open
		to author the title and description in.

	--no-edit
		Use the message from the first commit on the branch as pull request title
		and description without opening a text editor.

	-F, --file <FILE>
		Read the pull request title and description from <FILE>. Pass "-" to read
		from standard input instead. See ''--message'' for the formatting rules.

	-e, --edit
		Open the pull request title and description in a text editor before
		submitting. This can be used in combination with ''--message'' or ''--file''.

	-i, --issue <ISSUE>
		Convert <ISSUE> (referenced by its number) to a pull request.

		You can only convert issues authored by you or that which you have admin
		rights over. In most workflows it is not necessary to convert issues to
		pull requests; you can simply reference the original issue in the body of
		the new pull request.

	-o, --browse
		Open the new pull request in a web browser.

	-c, --copy
		Put the URL of the new pull request to clipboard instead of printing it.

	-p, --push
		Push the current branch to <HEAD> before creating the pull request.

	-b, --base <BASE>
		The base branch in the "[<OWNER>:]<BRANCH>" format. Defaults to the default
		branch of the upstream repository (usually "master").

		See the "CONVENTIONS" section of hub(1) for more information on how hub
		selects the defaults in case of multiple git remotes.

	-h, --head <HEAD>
		The head branch in "[<OWNER>:]<BRANCH>" format. Defaults to the currently
		checked out branch.

	-r, --reviewer <USERS>
		A comma-separated list (no spaces around the comma) of GitHub handles to
		request a review from.

	-a, --assign <USERS>
		A comma-separated list (no spaces around the comma) of GitHub handles to
		assign to this pull request.

	-M, --milestone <NAME>
		The milestone name to add to this pull request. Passing the milestone number
		is deprecated.

	-l, --labels <LABELS>
		A comma-separated list (no spaces around the comma) of labels to add to
		this pull request. Labels will be created if they do not already exist.

	-d, --draft
		Create the pull request as a draft.

	--no-maintainer-edits
		When creating a pull request from a fork, this disallows projects
		maintainers from being able to push to the head branch of this fork.
		Maintainer edits are allowed by default.

## Examples:
		$ hub pull-request
		[ opens a text editor for writing title and message ]
		[ creates a pull request for the current branch ]

		$ hub pull-request --base OWNER:master --head MYUSER:my-branch
		[ creates a pull request with explicit base and head branches ]

		$ hub pull-request --browse -m "My title"
		[ creates a pull request with the given title and opens it in a browser ]

		$ hub pull-request -F - --edit < path/to/message-template.md
		[ further edit the title and message received on standard input ]

## Configuration:

	* ''HUB_RETRY_TIMEOUT'':
		The maximum time to keep retrying after HTTP 422 on ''--push'' (default: 9).

## See also:

hub(1), hub-merge(1), hub-checkout(1)
`,
}

func init() {
	CmdRunner.Use(cmdPullRequest)
}

func pullRequest(cmd *Command, args *Args) {
	localRepo, err := github.LocalRepo()
	utils.Check(err)

	currentBranch, currentBranchErr := localRepo.CurrentBranch()

	baseProject, err := localRepo.MainProject()
	utils.Check(err)

	host, err := github.CurrentConfig().PromptForHost(baseProject.Host)
	if err != nil {
		utils.Check(github.FormatError("creating pull request", err))
	}
	client := github.NewClientWithHost(host)

	trackedBranch, headProject, _ := localRepo.RemoteBranchAndProject(host.User, false)
	if headProject == nil {
		utils.Check(fmt.Errorf("could not determine project for head branch"))
	}

	var (
		base, head string
	)

	if flagPullRequestBase := args.Flag.Value("--base"); flagPullRequestBase != "" {
		baseProject, base = parsePullRequestProject(baseProject, flagPullRequestBase)
	}

	if flagPullRequestHead := args.Flag.Value("--head"); flagPullRequestHead != "" {
		headProject, head = parsePullRequestProject(headProject, flagPullRequestHead)
	}

	baseRemote, _ := localRepo.RemoteForProject(baseProject)
	if base == "" && baseRemote != nil {
		base = localRepo.DefaultBranch(baseRemote).ShortName()
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

	force := args.Flag.Bool("--force")
	flagPullRequestPush := args.Flag.Bool("--push")

	if head == "" {
		if trackedBranch == nil {
			utils.Check(currentBranchErr)
			if !force && !flagPullRequestPush {
				branchRemote, branchMerge, err := branchTrackingInformation(currentBranch)
				if err != nil || (baseRemote != nil && branchRemote == baseRemote.Name && branchMerge.ShortName() == base) {
					if localRepo.RemoteForBranch(currentBranch, host.User) == nil {
						err = fmt.Errorf("Aborted: the current branch seems not yet pushed to a remote")
						err = fmt.Errorf("%s\n(use `-p` to push the branch or `-f` to skip this check)", err)
						utils.Check(err)
					}
				}
			}
			head = currentBranch.ShortName()
		} else {
			head = trackedBranch.ShortName()
		}
	}

	if headRepo, err := client.Repository(headProject); err == nil {
		headProject.Owner = headRepo.Owner.Login
		headProject.Name = headRepo.Name
	}

	fullBase := fmt.Sprintf("%s:%s", baseProject.Owner, base)
	fullHead := fmt.Sprintf("%s:%s", headProject.Owner, head)

	if !force && trackedBranch != nil {
		remoteCommits, err := git.RefList(trackedBranch.LongName(), "")
		if err == nil && len(remoteCommits) > 0 {
			err = fmt.Errorf("Aborted: %d commits are not yet pushed to %s", len(remoteCommits), trackedBranch.LongName())
			err = fmt.Errorf("%s\n(use `-f` to force submit a pull request anyway)", err)
			utils.Check(err)
		}
	}

	messageBuilder := &github.MessageBuilder{
		Filename: "PULLREQ_EDITMSG",
		Title:    "pull request",
	}

	baseTracking := base
	headTracking := head

	remote := baseRemote
	if remote != nil {
		baseTracking = fmt.Sprintf("%s/%s", remote.Name, base)
	}
	if remote == nil || !baseProject.SameAs(headProject) {
		remote, _ = localRepo.RemoteForProject(headProject)
	}
	if remote != nil {
		headTracking = fmt.Sprintf("%s/%s", remote.Name, head)
	}

	if flagPullRequestPush && remote == nil {
		utils.Check(fmt.Errorf("Can't find remote for %s", head))
	}

	messageBuilder.AddCommentedSection(fmt.Sprintf(`Requesting a pull to %s from %s

Write a message for this pull request. The first block
of text is the title and the rest is the description.`, fullBase, fullHead))

	flagPullRequestMessage := args.Flag.AllValues("--message")
	flagPullRequestEdit := args.Flag.Bool("--edit")
	flagPullRequestIssue := args.Flag.Value("--issue")
	if !args.Flag.HasReceived("--issue") && args.ParamsSize() > 0 {
		flagPullRequestIssue = parsePullRequestIssueNumber(args.GetParam(0))
	}

	if len(flagPullRequestMessage) > 0 {
		messageBuilder.Message = strings.Join(flagPullRequestMessage, "\n\n")
		messageBuilder.Edit = flagPullRequestEdit
	} else if args.Flag.HasReceived("--file") {
		messageBuilder.Message, err = msgFromFile(args.Flag.Value("--file"))
		utils.Check(err)
		messageBuilder.Edit = flagPullRequestEdit
	} else if args.Flag.Bool("--no-edit") {
		commits, _ := git.RefList(baseTracking, head)
		if len(commits) == 0 {
			utils.Check(fmt.Errorf("Aborted: no commits detected between %s and %s", baseTracking, head))
		}
		message, err := git.Show(commits[len(commits)-1])
		utils.Check(err)
		messageBuilder.Message = message
	} else if flagPullRequestIssue == "" {
		messageBuilder.Edit = true

		headForMessage := headTracking
		if flagPullRequestPush {
			headForMessage = head
		}

		message := ""

		commits, _ := git.RefList(baseTracking, headForMessage)
		if len(commits) == 1 {
			message, err = git.Show(commits[0])
			utils.Check(err)

			re := regexp.MustCompile(`\n(Co-authored-by|Signed-off-by):[^\n]+`)
			message = re.ReplaceAllString(message, "")
		} else if len(commits) > 1 {
			commitLogs, err := git.Log(baseTracking, headForMessage)
			utils.Check(err)

			if commitLogs != "" {
				messageBuilder.AddCommentedSection("\nChanges:\n\n" + strings.TrimSpace(commitLogs))
			}
		}

		workdir, _ := git.WorkdirName()
		if workdir != "" {
			template, _ := github.ReadTemplate(github.PullRequestTemplate, workdir)
			if template != "" {
				message = message + "\n\n\n" + template
			}
		}

		messageBuilder.Message = message
	}

	title, body, err := messageBuilder.Extract()
	utils.Check(err)

	if title == "" && flagPullRequestIssue == "" {
		utils.Check(fmt.Errorf("Aborting due to empty pull request title"))
	}

	if flagPullRequestPush {
		if args.Noop {
			args.Before(fmt.Sprintf("Would push to %s/%s", remote.Name, head), "")
		} else {
			err = git.Spawn("push", "--set-upstream", remote.Name, fmt.Sprintf("HEAD:%s", head))
			utils.Check(err)
		}
	}

	milestoneNumber, err := milestoneValueToNumber(args.Flag.Value("--milestone"), client, baseProject)
	utils.Check(err)

	var pullRequestURL string
	if args.Noop {
		args.Before(fmt.Sprintf("Would request a pull request to %s from %s", fullBase, fullHead), "")
		pullRequestURL = "PULL_REQUEST_URL"
	} else {
		params := map[string]interface{}{
			"base":                  base,
			"head":                  fullHead,
			"maintainer_can_modify": !args.Flag.Bool("--no-maintainer-edits"),
		}

		if args.Flag.Bool("--draft") {
			params["draft"] = true
		}

		if title != "" {
			params["title"] = title
			if body != "" {
				params["body"] = body
			}
		} else {
			issueNum, _ := strconv.Atoi(flagPullRequestIssue)
			params["issue"] = issueNum
		}

		startedAt := time.Now()
		numRetries := 0
		retryDelay := 2
		retryAllowance := 0
		if flagPullRequestPush {
			if allowanceFromEnv := os.Getenv("HUB_RETRY_TIMEOUT"); allowanceFromEnv != "" {
				retryAllowance, err = strconv.Atoi(allowanceFromEnv)
				utils.Check(err)
			} else {
				retryAllowance = 9
			}
		}

		var pr *github.PullRequest
		for {
			pr, err = client.CreatePullRequest(baseProject, params)
			if err != nil && strings.Contains(err.Error(), `Invalid value for "head"`) {
				if retryAllowance > 0 {
					retryAllowance -= retryDelay
					time.Sleep(time.Duration(retryDelay) * time.Second)
					retryDelay++
					numRetries++
				} else {
					if numRetries > 0 {
						duration := time.Since(startedAt)
						err = fmt.Errorf("%s\nGiven up after retrying for %.1f seconds.", err, duration.Seconds())
					}
					break
				}
			} else {
				break
			}
		}

		if err == nil {
			defer messageBuilder.Cleanup()
		}

		utils.Check(err)

		pullRequestURL = pr.HTMLURL

		params = map[string]interface{}{}
		flagPullRequestLabels := commaSeparated(args.Flag.AllValues("--labels"))
		if len(flagPullRequestLabels) > 0 {
			params["labels"] = flagPullRequestLabels
		}
		flagPullRequestAssignees := commaSeparated(args.Flag.AllValues("--assign"))
		if len(flagPullRequestAssignees) > 0 {
			params["assignees"] = flagPullRequestAssignees
		}
		if milestoneNumber > 0 {
			params["milestone"] = milestoneNumber
		}

		if len(params) > 0 {
			err = client.UpdateIssue(baseProject, pr.Number, params)
			utils.Check(err)
		}

		flagPullRequestReviewers := commaSeparated(args.Flag.AllValues("--reviewer"))
		if len(flagPullRequestReviewers) > 0 {
			userReviewers := []string{}
			teamReviewers := []string{}
			for _, reviewer := range flagPullRequestReviewers {
				if strings.Contains(reviewer, "/") {
					teamName := strings.SplitN(reviewer, "/", 2)[1]
					if !pr.HasRequestedTeam(teamName) {
						teamReviewers = append(teamReviewers, teamName)
					}
				} else if !pr.HasRequestedReviewer(reviewer) {
					userReviewers = append(userReviewers, reviewer)
				}
			}
			if len(userReviewers) > 0 || len(teamReviewers) > 0 {
				err = client.RequestReview(baseProject, pr.Number, map[string]interface{}{
					"reviewers":      userReviewers,
					"team_reviewers": teamReviewers,
				})
				utils.Check(err)
			}
		}
	}

	args.NoForward()
	printBrowseOrCopy(args, pullRequestURL, args.Flag.Bool("--browse"), args.Flag.Bool("--copy"))
}

func parsePullRequestProject(context *github.Project, s string) (p *github.Project, ref string) {
	p = context
	ref = s

	if strings.Contains(s, ":") {
		split := strings.SplitN(s, ":", 2)
		ref = split[1]
		var name string
		if !strings.Contains(split[0], "/") {
			name = context.Name
		}
		p = github.NewProject(split[0], name, context.Host)
	}

	return
}

func parsePullRequestIssueNumber(url string) string {
	u, e := github.ParseURL(url)
	if e != nil {
		return ""
	}

	r := regexp.MustCompile(`^issues\/(\d+)`)
	p := u.ProjectPath()
	if r.MatchString(p) {
		return r.FindStringSubmatch(p)[1]
	}

	return ""
}

func commaSeparated(l []string) []string {
	res := []string{}
	for _, i := range l {
		if i == "" {
			continue
		}
		res = append(res, strings.Split(i, ",")...)
	}
	return res
}
