package commands

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/github/hub/git"
	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var cmdPullRequest = &Command{
	Run: pullRequest,
	Usage: `
pull-request [-focp] [-b <BASE>] [-h <HEAD>] [-r <REVIEWERS> ] [-a <ASSIGNEES>] [-M <MILESTONE>] [-l <LABELS>]
pull-request -m <MESSAGE>
pull-request -F <FILE> [--edit]
pull-request -i <ISSUE>
`,
	Long: `Create a GitHub pull request.

## Options:
	-f, --force
		Skip the check for unpushed commits.

	--no-edit
		Use the commit message from the first commit on the branch without opening
		a text editor.

	-m, --message <MESSAGE>
		Use the first line of <MESSAGE> as pull request title, and the rest as pull
		request description.

	-F, --file <FILE>
		Read the pull request title and description from <FILE>.

	-e, --edit
		Further edit the contents of <FILE> in a text editor before submitting.

	-i, --issue <ISSUE>, <ISSUE-URL>
		(Deprecated) Convert <ISSUE> to a pull request.

	-o, --browse
		Open the new pull request in a web browser.

	-c, --copy
		Put the URL of the new pull request to clipboard instead of printing it.

	-p, --push
		Push the current branch to <HEAD> before creating the pull request.

	-b, --base <BASE>
		The base branch in "[OWNER:]BRANCH" format. Defaults to the default branch
		(usually "master").

	-h, --head <HEAD>
		The head branch in "[OWNER:]BRANCH" format. Defaults to the current branch.

	-r, --reviewer <USERS>
		A comma-separated list of GitHub handles to request a review from.

	-a, --assign <USERS>
		A comma-separated list of GitHub handles to assign to this pull request.

	-M, --milestone <NAME>
		The milestone name to add to this pull request. Passing the milestone number
		is deprecated.

	-l, --labels <LABELS>
		Add a comma-separated list of labels to this pull request.

## Configuration:

	HUB_RETRY_TIMEOUT=<SECONDS>
		The maximum time to keep retrying after HTTP 422 on '--push' (default: 9).

## See also:

hub(1), hub-merge(1), hub-checkout(1)
`,
}

var (
	flagPullRequestBase,
	flagPullRequestHead,
	flagPullRequestIssue,
	flagPullRequestMessage,
	flagPullRequestMilestone,
	flagPullRequestFile string

	flagPullRequestBrowse,
	flagPullRequestCopy,
	flagPullRequestEdit,
	flagPullRequestPush,
	flagPullRequestForce,
	flagPullRequestNoEdit bool

	flagPullRequestAssignees,
	flagPullRequestReviewers,
	flagPullRequestLabels listFlag
)

func init() {
	cmdPullRequest.Flag.StringVarP(&flagPullRequestBase, "base", "b", "", "BASE")
	cmdPullRequest.Flag.StringVarP(&flagPullRequestHead, "head", "h", "", "HEAD")
	cmdPullRequest.Flag.StringVarP(&flagPullRequestIssue, "issue", "i", "", "ISSUE")
	cmdPullRequest.Flag.BoolVarP(&flagPullRequestBrowse, "browse", "o", false, "BROWSE")
	cmdPullRequest.Flag.BoolVarP(&flagPullRequestCopy, "copy", "c", false, "COPY")
	cmdPullRequest.Flag.StringVarP(&flagPullRequestMessage, "message", "m", "", "MESSAGE")
	cmdPullRequest.Flag.BoolVarP(&flagPullRequestEdit, "edit", "e", false, "EDIT")
	cmdPullRequest.Flag.BoolVarP(&flagPullRequestPush, "push", "p", false, "PUSH")
	cmdPullRequest.Flag.BoolVarP(&flagPullRequestForce, "force", "f", false, "FORCE")
	cmdPullRequest.Flag.BoolVarP(&flagPullRequestNoEdit, "no-edit", "", false, "NO-EDIT")
	cmdPullRequest.Flag.StringVarP(&flagPullRequestFile, "file", "F", "", "FILE")
	cmdPullRequest.Flag.VarP(&flagPullRequestAssignees, "assign", "a", "USERS")
	cmdPullRequest.Flag.VarP(&flagPullRequestReviewers, "reviewer", "r", "USERS")
	cmdPullRequest.Flag.StringVarP(&flagPullRequestMilestone, "milestone", "M", "", "MILESTONE")
	cmdPullRequest.Flag.VarP(&flagPullRequestLabels, "labels", "l", "LABELS")

	CmdRunner.Use(cmdPullRequest)
}

func pullRequest(cmd *Command, args *Args) {
	localRepo, err := github.LocalRepo()
	utils.Check(err)

	currentBranch, err := localRepo.CurrentBranch()
	utils.Check(err)

	baseProject, err := localRepo.MainProject()
	utils.Check(err)

	host, err := github.CurrentConfig().PromptForHost(baseProject.Host)
	if err != nil {
		utils.Check(github.FormatError("creating pull request", err))
	}
	client := github.NewClientWithHost(host)

	trackedBranch, headProject, err := localRepo.RemoteBranchAndProject(host.User, false)
	utils.Check(err)

	var (
		base, head string
		force      bool
	)

	force = flagPullRequestForce

	if flagPullRequestBase != "" {
		baseProject, base = parsePullRequestProject(baseProject, flagPullRequestBase)
	}

	if flagPullRequestHead != "" {
		headProject, head = parsePullRequestProject(headProject, flagPullRequestHead)
	}

	if args.ParamsSize() == 1 {
		arg := args.RemoveParam(0)
		flagPullRequestIssue = parsePullRequestIssueNumber(arg)
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
		if trackedBranch == nil {
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
		remoteCommits, _ := git.RefList(trackedBranch.LongName(), "")
		if len(remoteCommits) > 0 {
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

	remote := gitRemoteForProject(baseProject)
	if remote != nil {
		baseTracking = fmt.Sprintf("%s/%s", remote.Name, base)
	}
	if remote == nil || !baseProject.SameAs(headProject) {
		remote = gitRemoteForProject(headProject)
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

	switch {
	case cmd.FlagPassed("message"):
		messageBuilder.Message = flagPullRequestMessage
		messageBuilder.Edit = flagPullRequestEdit
	case cmd.FlagPassed("file"):
		messageBuilder.Message, err = msgFromFile(flagPullRequestFile)
		utils.Check(err)
		messageBuilder.Edit = flagPullRequestEdit
	case flagPullRequestIssue == "":
		messageBuilder.Edit = true

		headForMessage := headTracking
		if flagPullRequestPush {
			headForMessage = head
		}

		message := ""
		commitLogs := ""

		commits, _ := git.RefList(baseTracking, headForMessage)
		if len(commits) == 1 {
			message, err = git.Show(commits[0])
			utils.Check(err)
		} else if len(commits) > 1 {
			commitLogs, err = git.Log(baseTracking, headForMessage)
			utils.Check(err)
		}

		if commitLogs != "" {
			messageBuilder.AddCommentedSection("\nChanges:\n\n" + strings.TrimSpace(commitLogs))
		}

		workdir, _ := git.WorkdirName()
		if workdir != "" {
			template, _ := github.ReadTemplate(github.PullRequestTemplate, workdir)
			if template != "" {
				message = message + "\n\n" + template
			}
		}

		messageBuilder.Message = message
	case flagPullRequestNoEdit:
		commits, _ := git.RefList(baseTracking, headTracking)
		if len(commits) == 0 {
			utils.Check(fmt.Errorf("No new commits on branch"))
		}
		message, err := git.Show(commits[len(commits)-1])
		utils.Check(err)
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

	milestoneNumber := 0
	if flagPullRequestMilestone != "" {
		// BC: Don't try to resolve milestone name if it's an integer
		milestoneNumber, err = strconv.Atoi(flagPullRequestMilestone)
		if err != nil {
			milestones, err := client.FetchMilestones(baseProject)
			utils.Check(err)
			milestoneNumber, err = findMilestoneNumber(milestones, flagPullRequestMilestone)
			utils.Check(err)
		}
	}

	var pullRequestURL string
	if args.Noop {
		args.Before(fmt.Sprintf("Would request a pull request to %s from %s", fullBase, fullHead), "")
		pullRequestURL = "PULL_REQUEST_URL"
	} else {
		params := map[string]interface{}{
			"base": base,
			"head": fullHead,
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
					retryDelay += 1
					numRetries += 1
				} else {
					if numRetries > 0 {
						duration := time.Now().Sub(startedAt)
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

		pullRequestURL = pr.HtmlUrl

		params = map[string]interface{}{}
		if len(flagPullRequestLabels) > 0 {
			params["labels"] = flagPullRequestLabels
		}
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

		if len(flagPullRequestReviewers) > 0 {
			userReviewers := []string{}
			teamReviewers := []string{}
			for _, reviewer := range flagPullRequestReviewers {
				if strings.Contains(reviewer, "/") {
					teamReviewers = append(teamReviewers, strings.SplitN(reviewer, "/", 2)[1])
				} else {
					userReviewers = append(userReviewers, reviewer)
				}
			}
			err = client.RequestReview(baseProject, pr.Number, map[string]interface{}{
				"reviewers":      userReviewers,
				"team_reviewers": teamReviewers,
			})
			utils.Check(err)
		}
	}

	if flagPullRequestIssue != "" {
		ui.Errorln("Warning: Issue to pull request conversion is deprecated and might not work in the future.")
	}

	args.NoForward()
	printBrowseOrCopy(args, pullRequestURL, flagPullRequestBrowse, flagPullRequestCopy)
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

func findMilestoneNumber(milestones []github.Milestone, name string) (int, error) {
	for _, milestone := range milestones {
		if strings.EqualFold(milestone.Title, name) {
			return milestone.Number, nil
		}
	}

	return 0, fmt.Errorf("error: no milestone found with name '%s'", name)
}
