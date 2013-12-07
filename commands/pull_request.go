package commands

import (
	"bufio"
	"fmt"
	"github.com/jingweno/gh/cmd"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
)

var cmdPullRequest = &Command{
	Run:   pullRequest,
	Usage: "pull-request [-f] [-i ISSUE] [-b BASE] [-h HEAD] [-m MESSAGE] [TITLE]",
	Short: "Open a pull request on GitHub",
	Long: `Opens a pull request on GitHub for the project that the "origin" remote
points to. The default head of the pull request is the current branch.
Both base and head of the pull request can be explicitly given in one of
the following formats: "branch", "owner:branch", "owner/repo:branch".
This command will abort operation if it detects that the current topic
branch has local commits that are not yet pushed to its upstream branch
on the remote. To skip this check, use -f.

If TITLE is omitted, a text editor will open in which title and body of
the pull request can be entered in the same manner as git commit message.

If instead of normal TITLE an issue number is given with -i, the pull
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
*/
func pullRequest(cmd *Command, args *Args) {
	localRepo := github.LocalRepo()

	currentBranch, err := localRepo.CurrentBranch()
	utils.Check(err)

	baseProject, err := localRepo.MainProject()
	utils.Check(err)

	headProject, err := localRepo.CurrentProject()
	utils.Check(err)

	gh := github.NewWithoutProject()
	gh.Project = baseProject

	var (
		base, head           string
		force, explicitOwner bool
	)
	force = flagPullRequestForce
	if flagPullRequestBase != "" {
		if strings.Contains(flagPullRequestBase, ":") {
			split := strings.SplitN(flagPullRequestBase, ":", 2)
			base = split[1]
			baseProject.Owner = split[0]
		} else {
			base = flagPullRequestBase
		}
	}

	if flagPullRequestHead != "" {
		if strings.Contains(flagPullRequestHead, ":") {
			split := strings.SplitN(flagPullRequestHead, ":", 2)
			head = split[1]
			headProject.Owner = split[0]
			explicitOwner = true
		} else {
			head = flagPullRequestHead
		}
	}

	if args.ParamsSize() == 1 {
		arg := args.RemoveParam(0)
		u, e := github.ParseURL(arg)
		r := regexp.MustCompile(`^issues\/(\d+)`)
		p := u.ProjectPath()
		if e == nil && r.MatchString(p) {
			flagPullRequestIssue = r.FindStringSubmatch(p)[1]
		}
	}

	if base == "" {
		masterBranch, err := localRepo.MasterBranch()
		utils.Check(err)
		base = masterBranch.ShortName()
	}

	trackedBranch, tberr := currentBranch.Upstream()
	if head == "" {
		if err == nil {
			if trackedBranch.IsRemote() {
				if reflect.DeepEqual(baseProject, headProject) && base == trackedBranch.ShortName() {
					e := fmt.Errorf(`Aborted: head branch is the same as base ("%s")`, base)
					e = fmt.Errorf("%s\n(use `-h <branch>` to specify an explicit pull request head)", e)
					utils.Check(e)
				}
			} else {
				// the current branch tracking another branch
				// pretend there's no upstream at all
				tberr = fmt.Errorf("No upstream found for current branch")
			}
		}

		if tberr == nil {
			head = trackedBranch.ShortName()
		} else {
			head = currentBranch.ShortName()
		}
	}

	// when no tracking, assume remote branch is published under active user's fork
	if tberr != nil && !explicitOwner && gh.Config.User != headProject.Owner {
		headProject = github.NewProjectFromString(headProject.Name)
	}

	var title, body string
	if flagPullRequestMessage != "" {
		title, body = readMsg(flagPullRequestMessage)
	}

	if flagPullRequestFile != "" {
		var (
			content []byte
			err     error
		)
		if flagPullRequestFile == "-" {
			content, err = ioutil.ReadAll(os.Stdin)
		} else {
			content, err = ioutil.ReadFile(flagPullRequestFile)
		}
		utils.Check(err)
		title, body = readMsg(string(content))
	}

	fullBase := fmt.Sprintf("%s:%s", baseProject.Owner, base)
	fullHead := fmt.Sprintf("%s:%s", headProject.Owner, head)

	commits, _ := git.RefList(base, head)
	if !force && tberr == nil && len(commits) > 0 {
		err = fmt.Errorf("Aborted: %d commits are not yet pushed to %s", len(commits), trackedBranch.LongName())
		err = fmt.Errorf("%s\n(use `-f` to force submit a pull request anyway)", err)
		utils.Check(err)
	}

	if title == "" && flagPullRequestIssue == "" {
		t, b, err := writePullRequestTitleAndBody(base, head, fullBase, fullHead, commits)
		utils.Check(err)
		title = t
		body = b
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
			pr, err := gh.CreatePullRequest(base, fullHead, title, body)
			utils.Check(err)
			pullRequestURL = pr.HTMLURL
		}

		if flagPullRequestIssue != "" {
			pr, err := gh.CreatePullRequestForIssue(base, fullHead, flagPullRequestIssue)
			utils.Check(err)
			pullRequestURL = pr.HTMLURL
		}
	}

	args.Replace("echo", "", pullRequestURL)
}

func writePullRequestTitleAndBody(base, head, fullBase, fullHead string, commits []string) (title, body string, err error) {
	messageFile, err := git.PullReqMsgFile()
	if err != nil {
		return
	}
	defer os.Remove(messageFile)

	err = writePullRequestChanges(base, head, fullBase, fullHead, commits, messageFile)
	if err != nil {
		return
	}

	editor, err := git.Editor()
	if err != nil {
		return
	}

	err = editTitleAndBody(editor, messageFile)
	if err != nil {
		err = fmt.Errorf("error using text editor for pull request message")
		return
	}

	title, body, err = readTitleAndBody(messageFile)
	if err != nil {
		return
	}

	return
}

func writePullRequestChanges(base, head, fullBase, fullHead string, commits []string, messageFile string) error {
	var defaultMsg, commitSummary string
	if len(commits) == 1 {
		defaultMsg, err := git.Show(commits[0])
		if err != nil {
			return err
		}
		defaultMsg = fmt.Sprintf("%s\n", defaultMsg)
	} else if len(commits) > 1 {
		commitLogs, err := git.Log(base, head)
		if err != nil {
			return err
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

	return ioutil.WriteFile(messageFile, []byte(message), 0644)
}

func editTitleAndBody(editor, messageFile string) error {
	editCmd := cmd.New(editor)
	r := regexp.MustCompile("[mg]?vi[m]$")
	if r.MatchString(editor) {
		editCmd.WithArg("-c")
		editCmd.WithArg("set ft=gitcommit tw=0 wrap lbr")
	}
	editCmd.WithArg(messageFile)

	return editCmd.Exec()
}

func readTitleAndBody(messageFile string) (title, body string, err error) {
	f, err := os.Open(messageFile)
	defer f.Close()
	if err != nil {
		return "", "", err
	}

	reader := bufio.NewReader(f)

	return readTitleAndBodyFrom(reader)
}

func readTitleAndBodyFrom(reader *bufio.Reader) (title, body string, err error) {
	r := regexp.MustCompile("\\S")
	var titleParts, bodyParts []string

	line, err := readLine(reader)
	for err == nil {
		if strings.HasPrefix(line, "#") {
			break
		}

		if len(bodyParts) == 0 && r.MatchString(line) {
			titleParts = append(titleParts, line)
		} else {
			bodyParts = append(bodyParts, line)
		}

		line, err = readLine(reader)
	}

	title = strings.Join(titleParts, " ")
	title = strings.TrimSpace(title)

	body = strings.Join(bodyParts, "\n")
	body = strings.TrimSpace(body)

	return title, body, nil
}

func readLine(r *bufio.Reader) (string, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}

	return string(ln), err
}

func readMsg(msg string) (title, body string) {
	split := strings.SplitN(msg, "\n\n", 2)
	title = strings.TrimSpace(split[0])
	if len(split) > 1 {
		body = strings.TrimSpace(split[1])
	}

	return
}
