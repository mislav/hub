package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/github/hub/git"
	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
	"github.com/octokit/go-octokit/octokit"
)

var (
	cmdPullRequest = &Command{
		Run: pullRequest,
		IgnoreUnknownSubcommands: true,
		Usage: `
pull-request [-fo] [-b <BASE>] [-h <HEAD>] [-a <USER>] [-M <MILESTONE>] [-l <LABELS>]
pull-request -m <MESSAGE>
pull-request -F <FILE>
pull-request -i <ISSUE>
pull-request info [-b <BASE>] [-h <HEAD>]
`,
		Long: `Manage GitHub pull requests.

## Commands:

With no subcommand, creates a GitHub pull request.

	* _info_:
		Show info for the GitHub pull request corresponding to a branch.

## Options:
	-f, --force
		Skip the check for unpushed commits.

	-m, --message <MESSAGE>
		Use the first line of <MESSAGE> as pull request title, and the rest as pull
		request description.

	-F, --file <FILE>
		Read the pull request title and description from <FILE>.

	-i, --issue <ISSUE>, <ISSUE-URL>
		(Deprecated) Convert <ISSUE> to a pull request.

	-o, --browse
		Open the new pull request in a web browser.

	-b, --base <BASE>
		The base branch in "[OWNER:]BRANCH" format. Defaults to the default branch
		(usually "master").

	-h, --head <HEAD>
		The base branch in "[OWNER:]BRANCH" format. Defaults to the current branch.

	-a, --assign <USER>
		Assign GitHub <USER> to this pull request.

	-M, --milestone <ID>
		Add this pull request to a GitHub milestone with id <ID>.

	-l, --labels <LABELS>
		Add a comma-separated list of labels to this pull request.

## See also:

hub(1), hub-merge(1), hub-checkout(1)
`,
	}
	cmdPullRequestInfo = &Command{
		Key:   "info",
		Run:   pullRequestInfo,
	}
)

var (
	flagPullRequestBase,
	flagPullRequestHead,
	flagPullRequestIssue,
	flagPullRequestMessage,
	flagPullRequestAssignee,
	flagPullRequestLabels,
	flagPullRequestFile,
	flagPullRequestInfoBase,
	flagPullRequestInfoHead string
	flagPullRequestBrowse,
	flagPullRequestForce bool
	flagPullRequestMilestone uint64
)

func init() {
	cmdPullRequest.Flag.StringVarP(&flagPullRequestBase, "base", "b", "", "BASE")
	cmdPullRequest.Flag.StringVarP(&flagPullRequestHead, "head", "h", "", "HEAD")
	cmdPullRequest.Flag.StringVarP(&flagPullRequestIssue, "issue", "i", "", "ISSUE")
	cmdPullRequest.Flag.BoolVarP(&flagPullRequestBrowse, "browse", "o", false, "BROWSE")
	cmdPullRequest.Flag.StringVarP(&flagPullRequestMessage, "message", "m", "", "MESSAGE")
	cmdPullRequest.Flag.BoolVarP(&flagPullRequestForce, "force", "f", false, "FORCE")
	cmdPullRequest.Flag.StringVarP(&flagPullRequestFile, "file", "F", "", "FILE")
	cmdPullRequest.Flag.StringVarP(&flagPullRequestAssignee, "assign", "a", "", "USER")
	cmdPullRequest.Flag.Uint64VarP(&flagPullRequestMilestone, "milestone", "M", 0, "MILESTONE")
	cmdPullRequest.Flag.StringVarP(&flagPullRequestLabels, "labels", "l", "", "LABELS")

	cmdPullRequestInfo.Flag.StringVarP(&flagPullRequestInfoBase, "base", "b", "", "BASE")
	cmdPullRequestInfo.Flag.StringVarP(&flagPullRequestInfoHead, "head", "h", "", "HEAD")

	cmdPullRequest.Use(cmdPullRequestInfo)

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

	title, body, err := getTitleAndBodyFromFlags(flagPullRequestMessage, flagPullRequestFile)
	utils.Check(err)

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
	if title == "" && flagPullRequestIssue == "" {
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

		message, err := pullRequestChangesMessage(baseTracking, headTracking, fullBase, fullHead)
		utils.Check(err)

		editor, err = github.NewEditor("PULLREQ", "pull request", message)
		utils.Check(err)

		title, body, err = editor.EditTitleAndBody()
		utils.Check(err)
	}

	if title == "" && flagPullRequestIssue == "" {
		utils.Check(fmt.Errorf("Aborting due to empty pull request title"))
	}

	var pullRequestURL string
	if args.Noop {
		args.Before(fmt.Sprintf("Would request a pull request to %s from %s", fullBase, fullHead), "")
		pullRequestURL = "PULL_REQUEST_URL"
	} else {
		var (
			pr  *octokit.PullRequest
			err error
		)

		client := github.NewClientWithHost(host)
		if title != "" {
			pr, err = client.CreatePullRequest(baseProject, base, fullHead, title, body)
		} else if flagPullRequestIssue != "" {
			pr, err = client.CreatePullRequestForIssue(baseProject, base, fullHead, flagPullRequestIssue)
		}

		if err == nil && editor != nil {
			defer editor.DeleteFile()
		}

		utils.Check(err)

		pullRequestURL = pr.HTMLURL

		if flagPullRequestAssignee != "" || flagPullRequestMilestone > 0 ||
			flagPullRequestLabels != "" {

			labels := []string{}
			for _, label := range strings.Split(flagPullRequestLabels, ",") {
				if label != "" {
					labels = append(labels, label)
				}
			}

			params := octokit.IssueParams{
				Assignee:  flagPullRequestAssignee,
				Milestone: flagPullRequestMilestone,
				Labels:    labels,
			}

			err = client.UpdateIssue(baseProject, pr.Number, params)
			utils.Check(err)
		}
	}

	if flagPullRequestBrowse {
		launcher, err := utils.BrowserLauncher()
		utils.Check(err)
		args.Replace(launcher[0], "", launcher[1:]...)
		args.AppendParams(pullRequestURL)
	} else {
		args.Replace("echo", "", pullRequestURL)
	}

	if flagPullRequestIssue != "" {
		args.After("echo", "Warning: Issue to pull request conversion is deprecated and might not work in the future.")
	}
}

func pullRequestChangesMessage(base, head, fullBase, fullHead string) (string, error) {
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

// pullRequestInfo fetches info for a pull request.
func pullRequestInfo(cmd *Command, args *Args) {
	localRepo, err := github.LocalRepo()
	utils.Check(err)

	currentBranch, err := localRepo.CurrentBranch()
	utils.Check(err)

	baseProject, err := localRepo.MainProject()
	utils.Check(err)

	host, err := github.CurrentConfig().PromptForHost(baseProject.Host)
	if err != nil {
		utils.Check(github.FormatError("getting pull request info", err))
	}

	trackedBranch, headProject, err := localRepo.RemoteBranchAndProject(host.User, false)

	var (
		base, head string
	)

	if flagPullRequestInfoBase != "" {
		baseProject, base = parsePullRequestProject(baseProject, flagPullRequestInfoBase)
	}

	if flagPullRequestInfoHead != "" {
		headProject, head = parsePullRequestProject(headProject, flagPullRequestInfoHead)
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

	fullHead := fmt.Sprintf("%s:%s", headProject.Owner, head)

	client := github.NewClientWithHost(host)
	pr, err := client.FindPullRequest(baseProject, base, fullHead)
	utils.Check(err)

	comments, err := client.PullRequestComments(pr)
	utils.Check(err)

	ui.Printf("Title: %q\n", trimToOneLine(pr.Title))
	for _, comment := range comments {
		ui.Printf("* %s - %s: %s\n", comment.UpdatedAt.Format("2006-01-02 15:04:05"), comment.User.Login, trimToOneLine(comment.Body))
	}

	args.Replace("echo", "", "URL:", pr.HTMLURL)
}

func trimToOneLine(s string) string {
	if i := strings.IndexAny(s, "\r\n"); i != -1 {
		s = s[:i] + "…"
	}
	if len(s) >= 80 {
		s = s[:80] + "…"
	}
	return s
}
