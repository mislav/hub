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
pull-request [-foc] [-b <BASE>] [-h <HEAD>] [-r <REVIEWERS> ] [-a <ASSIGNEES>] [-M <MILESTONE>] [-l <LABELS>]
pull-request -m <MESSAGE>
pull-request -F <FILE> [--edit]
pull-request -i <ISSUE>
`,
	Long: `Create a GitHub pull request.

## Options:
	-f, --force
		Skip the check for unpushed commits.

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

	-M, --milestone <ID>
		Add this pull request to a GitHub milestone with id <ID>.

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
	flagPullRequestFile string

	flagPullRequestBrowse,
	flagPullRequestCopy,
	flagPullRequestEdit,
	flagPullRequestPush,
	flagPullRequestForce bool

	flagPullRequestMilestone uint64

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
	cmdPullRequest.Flag.StringVarP(&flagPullRequestFile, "file", "F", "", "FILE")
	cmdPullRequest.Flag.VarP(&flagPullRequestAssignees, "assign", "a", "USERS")
	cmdPullRequest.Flag.VarP(&flagPullRequestReviewers, "reviewer", "r", "USERS")
	cmdPullRequest.Flag.Uint64VarP(&flagPullRequestMilestone, "milestone", "M", 0, "MILESTONE")
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

	var editor *github.Editor
	var title, body string

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

	if cmd.FlagPassed("message") {
		title, body = readMsg(flagPullRequestMessage)
	} else if cmd.FlagPassed("file") {
		title, body, editor, err = readMsgFromFile(flagPullRequestFile, flagPullRequestEdit, "PULLREQ", "pull request")
		utils.Check(err)
	} else if flagPullRequestIssue == "" {
		headForMessage := headTracking
		if flagPullRequestPush {
			headForMessage = head
		}

		message, err := createPullRequestMessage(baseTracking, headForMessage, fullBase, fullHead)
		utils.Check(err)

		editor, err = github.NewEditor("PULLREQ", "pull request", message)
		utils.Check(err)

		title, body, err = editor.EditTitleAndBody()
		utils.Check(err)
	}

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

		if err == nil && editor != nil {
			defer editor.DeleteFile()
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
		if flagPullRequestMilestone > 0 {
			params["milestone"] = flagPullRequestMilestone
		}

		if len(params) > 0 {
			err = client.UpdateIssue(baseProject, pr.Number, params)
			utils.Check(err)
		}

		if len(flagPullRequestReviewers) > 0 {
			err = client.RequestReview(baseProject, pr.Number, map[string]interface{}{"reviewers": flagPullRequestReviewers})
			utils.Check(err)
		}
	}

	if flagPullRequestIssue != "" {
		ui.Errorln("Warning: Issue to pull request conversion is deprecated and might not work in the future.")
	}

	args.NoForward()
	printBrowseOrCopy(args, pullRequestURL, flagPullRequestBrowse, flagPullRequestCopy)
}

func createPullRequestMessage(base, head, fullBase, fullHead string) (string, error) {
	var (
		defaultMsg string
		commitLogs string
		err        error
	)

	commits, _ := git.RefList(base, head)
	if len(commits) == 1 {
		defaultMsg, err = git.Show(commits[0])
		if err != nil {
			return "", err
		}
	} else if len(commits) > 1 {
		commitLogs, err = git.Log(base, head)
		if err != nil {
			return "", err
		}
	}

	workdir, _ := git.WorkdirName()
	if workdir != "" {
		template, err := github.ReadTemplate(github.PullRequestTemplate, workdir)
		if err != nil {
			return "", err
		} else if template != "" {
			if defaultMsg == "" {
				defaultMsg = "\n\n" + template
			} else {
				parts := strings.SplitN(defaultMsg, "\n\n", 2)
				defaultMsg = parts[0] + "\n\n" + template
				if len(parts) > 1 && parts[1] != "" {
					defaultMsg = defaultMsg + "\n\n" + parts[1]
				}
			}
		}
	}

	cs := git.CommentChar()

	return renderPullRequestTpl(defaultMsg, cs, fullBase, fullHead, commitLogs)
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
