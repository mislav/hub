package commands

import (
	"fmt"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"reflect"
	"regexp"
	"strings"
)

var cmdPullRequest = &Command{
	Run:   pullRequest,
	Usage: "pull-request [-f] [-m <MESSAGE>|-F <FILE>|-i <ISSUE>|<ISSUE-URL>] [-b <BASE>] [-h <HEAD>] ",
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
	flagPullRequestForce bool
)

func init() {
	cmdPullRequest.Flag.StringVar(&flagPullRequestBase, "b", "", "BASE")
	cmdPullRequest.Flag.StringVar(&flagPullRequestHead, "h", "", "HEAD")
	cmdPullRequest.Flag.StringVar(&flagPullRequestIssue, "i", "", "ISSUE")
	cmdPullRequest.Flag.StringVar(&flagPullRequestMessage, "m", "", "MESSAGE")
	cmdPullRequest.Flag.BoolVar(&flagPullRequestForce, "f", false, "FORCE")
	cmdPullRequest.Flag.StringVar(&flagPullRequestFile, "F", "", "FILE")
	cmdPullRequest.Flag.StringVar(&flagPullRequestFile, "file", "", "FILE")

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
	localRepo := github.LocalRepo()

	currentBranch, err := localRepo.CurrentBranch()
	utils.Check(err)

	baseProject, err := localRepo.MainProject()
	utils.Check(err)

	client := github.NewClient(baseProject.Host)

	trackedBranch, headProject, err := localRepo.RemoteBranchAndProject(client.Credentials.User)
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

	if head == "" {
		if !trackedBranch.IsRemote() {
			// the current branch tracking another branch
			// pretend there's no upstream at all
			trackedBranch = nil
		} else {
			if reflect.DeepEqual(baseProject, headProject) && base == trackedBranch.ShortName() {
				e := fmt.Errorf(`Aborted: head branch is the same as base ("%s")`, base)
				e = fmt.Errorf("%s\n(use `-h <branch>` to specify an explicit pull request head)", e)
				utils.Check(e)
			}
		}

		if trackedBranch == nil {
			head = currentBranch.ShortName()
		} else {
			head = trackedBranch.ShortName()
		}
	}

	title, body, err := github.GetTitleAndBodyFromFlags(flagPullRequestMessage, flagPullRequestFile)
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

	if title == "" && flagPullRequestIssue == "" {
		commits, _ := git.RefList(base, head)
		title, body, err = writePullRequestTitleAndBody(base, head, fullBase, fullHead, commits)
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
		if title != "" {
			pr, err := client.CreatePullRequest(baseProject, base, fullHead, title, body)
			utils.Check(err)
			pullRequestURL = pr.HTMLURL
		}

		if flagPullRequestIssue != "" {
			pr, err := client.CreatePullRequestForIssue(baseProject, base, fullHead, flagPullRequestIssue)
			utils.Check(err)
			pullRequestURL = pr.HTMLURL
		}
	}

	args.Replace("echo", "", pullRequestURL)
	if flagPullRequestIssue != "" {
		args.After("echo", "Warning: Issue to pull request conversion is deprecated and might not work in the future.")
	}
}

func writePullRequestTitleAndBody(base, head, fullBase, fullHead string, commits []string) (title, body string, err error) {
	message, err := pullRequestChangesMessage(base, head, fullBase, fullHead, commits)
	utils.Check(err)

	return github.GetTitleAndBodyFromEditor("PULLREQ", message)
}

func pullRequestChangesMessage(base, head, fullBase, fullHead string, commits []string) (string, error) {
	var defaultMsg, commitSummary string
	if len(commits) == 1 {
		msg, err := git.Show(commits[0])
		if err != nil {
			return "", err
		}
		defaultMsg = fmt.Sprintf("%s\n", msg)
	} else if len(commits) > 1 {
		commitLogs, err := git.Log(base, head)
		if err != nil {
			return "", err
		}

		if len(commitLogs) > 0 {
			startRegexp := regexp.MustCompilePOSIX("^")
			endRegexp := regexp.MustCompilePOSIX(" +$")

			commitLogs = strings.TrimSpace(commitLogs)
			commitLogs = startRegexp.ReplaceAllString(commitLogs, "# ")
			commitLogs = endRegexp.ReplaceAllString(commitLogs, "")
			commitSummary = `
#
# Changes:
#
%s`
			commitSummary = fmt.Sprintf(commitSummary, commitLogs)
		}
	}

	message := `%s
# Requesting a pull to %s from %s
#
# Write a message for this pull request. The first block
# of the text is the title and the rest is description.%s
`
	message = fmt.Sprintf(message, defaultMsg, fullBase, fullHead, commitSummary)

	return message, nil
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
