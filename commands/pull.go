package commands

import (
	"bufio"
	"fmt"
	"github.com/jingweno/gh/cmd"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var cmdPull = &Command{
	Run:   pull,
	Usage: "pull [-f] [TITLE|-i ISSUE] [-b BASE] [-h HEAD]",
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
	cmdPull.Flag.StringVar(&flagPullRequestBase, "b", "master", "BASE")
	cmdPull.Flag.StringVar(&flagPullRequestHead, "h", "", "HEAD")
	cmdPull.Flag.StringVar(&flagPullRequestIssue, "i", "", "ISSUE")
}

func pull(cmd *Command, args []string) {
	var title, body string
	if len(args) == 1 {
		title = args[0]
	}

	gh := github.New()
	repo := gh.Project.LocalRepoWith(flagPullRequestBase, flagPullRequestHead)
	if title == "" && flagPullRequestIssue == "" {
		messageFile, err := git.PullReqMsgFile()
		utils.Check(err)

		err = writePullRequestChanges(repo, messageFile)
		utils.Check(err)

		editorPath, err := git.EditorPath()
		utils.Check(err)

		err = editTitleAndBody(editorPath, messageFile)
		utils.Check(err)

		title, body, err = readTitleAndBody(messageFile)
		utils.Check(err)
	}

	if title == "" && flagPullRequestIssue == "" {
		log.Fatal("Aborting due to empty pull request title")
	}

	var pullRequestURL string
	var err error
	if title != "" {
		pullRequestURL, err = gh.CreatePullRequest(repo.Base, repo.Head, title, body)
	}
	if flagPullRequestIssue != "" {
		pullRequestURL, err = gh.CreatePullRequestForIssue(repo.Base, repo.Head, flagPullRequestIssue)
	}

	utils.Check(err)

	fmt.Println(pullRequestURL)
}

func writePullRequestChanges(repo *github.Repo, messageFile string) error {
	message := `
# Requesting a pull to %s from %s
#
# Write a message for this pull reuqest. The first block
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
