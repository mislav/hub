package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/github/hub/Godeps/_workspace/src/github.com/octokit/go-octokit/octokit"
	"github.com/github/hub/git"
	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var cmdPullRequest = &Command{
	Run:   pullRequest,
	Usage: "pull-request [-f] [-m <MESSAGE>|-F <FILE>|-i <ISSUE>|<ISSUE-URL>] [-o] [-b <BASE>] [-h <HEAD>] ",
	Short: "Open a pull request on GitHub",
	Long: `Opens a pull request on GitHub for the project that the "origin" remote
points to. The default head of the pull request is the current branch.
Both base and head of the pull request can be explicitly given in one of
the following formats: "branch", "owner:branch", "owner/repo:branch".
This command will abort operation if it detects that the current topic
branch has local commits that are not yet pushed to its upstream branch
on the remote. To skip this check, use "-f".

Without <MESSAGE> or <FILE>, a text editor will open in which title and body
of the pull request can be entered in the same manner as git commit message.
Pull request message can also be passed via stdin with "-F -".

If instead of normal <TITLE> an issue number is given with "-i", the pull
request will be attached to an existing GitHub issue. Alternatively, instead
of title you can paste a full URL to an issue on GitHub.
`,
}

var (
	flagPullRequestBase,
	flagPullRequestHead,
	flagPullRequestIssue,
	flagPullRequestMessage,
	flagPullRequestFile string
	flagPullRequestBrowse,
	flagPullRequestForce bool
)

func init() {
	cmdPullRequest.Flag.StringVarP(&flagPullRequestBase, "base", "b", "", "BASE")
	cmdPullRequest.Flag.StringVarP(&flagPullRequestHead, "head", "h", "", "HEAD")
	cmdPullRequest.Flag.StringVarP(&flagPullRequestIssue, "issue", "i", "", "ISSUE")
	cmdPullRequest.Flag.BoolVarP(&flagPullRequestBrowse, "browse", "o", false, "BROWSE")
	cmdPullRequest.Flag.StringVarP(&flagPullRequestMessage, "message", "m", "", "MESSAGE")
	cmdPullRequest.Flag.BoolVarP(&flagPullRequestForce, "force", "f", false, "FORCE")
	cmdPullRequest.Flag.StringVarP(&flagPullRequestFile, "file", "F", "", "FILE")

	CmdRunner.Use(cmdPullRequest)
}

/*
  # while on a topic branch called "feature":
  $ gh pull-request
  [ opens text editor to edit title & body for the request ]
  [ opened pull request on GitHub for "YOUR_USER:feature" ]

  # explicit pull base & head:
  $ gh pull-request -b jingweno:master -h jingweno:feature

  $ gh pull-request -m "title\n\nbody"
  [ create pull request with title & body  ]

  $ gh pull-request -i 123
  [ attached pull request to issue #123 ]

  $ gh pull-request https://github.com/jingweno/gh/pull/123
  [ attached pull request to issue #123 ]

  $ gh pull-request -F FILE
  [ create pull request with title & body from FILE ]
*/
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
