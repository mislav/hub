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
	"regexp"
	"strings"
)

var cmdPullRequest = &Command{
	Run:   pullRequest,
	Usage: "pull-request [-f] [-i ISSUE] [-b BASE] [-d HEAD] [TITLE]",
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

var flagPullRequestBase, flagPullRequestHead, flagPullRequestIssue string

func init() {
	cmdPullRequest.Flag.StringVar(&flagPullRequestBase, "b", "master", "BASE")
	cmdPullRequest.Flag.StringVar(&flagPullRequestHead, "d", "", "HEAD")
	cmdPullRequest.Flag.StringVar(&flagPullRequestIssue, "i", "", "ISSUE")
}

/*
  # while on a topic branch called "feature":
  $ gh pull-request
  [ opens text editor to edit title & body for the request ]
  [ opened pull request on GitHub for "YOUR_USER:feature" ]

  # explicit pull base & head:
  $ gh pull-request -b jingweno:master -h jingweno:feature

  $ gh pull-request -i 123
  [ attached pull request to issue #123 ]
*/
func pullRequest(cmd *Command, args *Args) {
	var (
		title, body string
		err         error
	)
	if args.ParamsSize() == 1 {
		title = args.RemoveParam(0)
	}

	gh := github.New()
	repo := gh.Project.LocalRepoWith(flagPullRequestBase, flagPullRequestHead)
	if title == "" && flagPullRequestIssue == "" {
		title, body, err = writePullRequestTitleAndBody(repo)
	}

	if title == "" && flagPullRequestIssue == "" {
		utils.Check(fmt.Errorf("Aborting due to empty pull request title"))
	}

	var pullRequestURL string
	if args.Noop {
		args.Before(fmt.Sprintf("Would request a pull request to %s from %s", repo.FullBase(), repo.FullHead()), "")
		pullRequestURL = "PULL_REQUEST_URL"
	} else {
		if title != "" {
			pullRequestURL, err = gh.CreatePullRequest(repo.Base, repo.Head, title, body)
		}
		utils.Check(err)

		if flagPullRequestIssue != "" {
			pullRequestURL, err = gh.CreatePullRequestForIssue(repo.Base, repo.Head, flagPullRequestIssue)
		}
		utils.Check(err)

	}

	args.Replace("echo", "", pullRequestURL)
}

func writePullRequestTitleAndBody(repo *github.Repo) (title, body string, err error) {
	messageFile, err := git.PullReqMsgFile()
	if err != nil {
		return
	}

	err = writePullRequestChanges(repo, messageFile)
	if err != nil {
		return
	}

	editorPath, err := git.EditorPath()
	if err != nil {
		return
	}

	err = editTitleAndBody(editorPath, messageFile)
	if err != nil {
		return
	}

	title, body, err = readTitleAndBody(messageFile)
	if err != nil {
		return
	}

	err = os.Remove(messageFile)

	return
}

func writePullRequestChanges(repo *github.Repo, messageFile string) error {
	message := `
# Requesting a pull to %s from %s
#
# Write a message for this pull request. The first block
# of the text is the title and the rest is description.%s
`
	startRegexp := regexp.MustCompilePOSIX("^")
	endRegexp := regexp.MustCompilePOSIX(" +$")

	commitLogs, _ := git.Log(repo.Base, repo.Head)
	var changesMsg string
	if len(commitLogs) > 0 {
		commitLogs = strings.TrimSpace(commitLogs)
		commitLogs = startRegexp.ReplaceAllString(commitLogs, "# ")
		commitLogs = endRegexp.ReplaceAllString(commitLogs, "")
		changesMsg = `
#
# Changes:
#
%s`
		changesMsg = fmt.Sprintf(changesMsg, commitLogs)
	}

	message = fmt.Sprintf(message, repo.FullBase(), repo.FullHead(), changesMsg)

	return ioutil.WriteFile(messageFile, []byte(message), 0644)
}

func editTitleAndBody(editorPath, messageFile string) error {
	editCmd := cmd.New(editorPath)
	r := regexp.MustCompile("[mg]?vi[m]$")
	if r.MatchString(editorPath) {
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

	line, err := readln(reader)
	for err == nil {
		if strings.HasPrefix(line, "#") {
			break
		}
		if len(bodyParts) == 0 && r.MatchString(line) {
			titleParts = append(titleParts, line)
		} else {
			bodyParts = append(bodyParts, line)
		}
		line, err = readln(reader)
	}

	title = strings.Join(titleParts, " ")
	title = strings.TrimSpace(title)

	body = strings.Join(bodyParts, "\n")
	body = strings.TrimSpace(body)

	return title, body, nil
}

func readln(r *bufio.Reader) (string, error) {
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
