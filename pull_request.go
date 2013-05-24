package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var cmdPullRequest = &Command{
	Run:   pullRequest,
	Usage: "pull-request [-f] [TITLE|-i ISSUE] [-b BASE] [-h HEAD]",
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

var flagPullRequestBase, flagPullRequestHead string

func init() {
	head, _ := FetchGitHead()

	cmdPullRequest.Flag.StringVar(&flagPullRequestBase, "b", "master", "BASE")
	cmdPullRequest.Flag.StringVar(&flagPullRequestHead, "h", head, "HEAD")
}

func pullRequest(cmd *Command, args []string) {
	repo := NewRepo()
	repo.Base = flagPullRequestBase
	repo.Head = flagPullRequestHead

	messageFile := filepath.Join(repo.Dir, "PULLREQ_EDITMSG")

	err := writePullRequestChanges(repo, messageFile)
	check(err)

	editCmd := buildEditCommand(repo, messageFile)
	err = editCmd.Exec()
	check(err)

	title, body, err := readTitleAndBodyFromFile(messageFile)
	check(err)

	if len(title) == 0 {
		log.Fatal("Aborting due to empty pull request title")
	}

	params := PullRequestParams{title, body, flagPullRequestBase, flagPullRequestHead}
	gh := NewGitHub()
	pullRequestResponse, err := gh.CreatePullRequest(repo.Owner, repo.Project, params)
	check(err)

	fmt.Println(pullRequestResponse.HtmlUrl)
}

func writePullRequestChanges(repo *Repo, messageFile string) error {
	message := `
# Requesting a pull to %s from %s
#
# Write a message for this pull reuqest. The first block
# of the text is the title and the rest is description.
#
# Changes:
#
%s
`
	startRegexp := regexp.MustCompilePOSIX("^")
	endRegexp := regexp.MustCompilePOSIX(" +$")

	commitLogs, err := FetchGitCommitLogs(repo.Base, repo.Head)
	check(err)

	commitLogs = strings.TrimSpace(commitLogs)
	commitLogs = startRegexp.ReplaceAllString(commitLogs, "# ")
	commitLogs = endRegexp.ReplaceAllString(commitLogs, "")

	message = fmt.Sprintf(message, repo.FullBase(), repo.FullHead(), commitLogs)

	return ioutil.WriteFile(messageFile, []byte(message), 0644)
}

func getLocalBranch(branchName string) string {
	result := strings.Split(branchName, ":")

	return result[len(result)-1]
}

func buildEditCommand(repo *Repo, messageFile string) *ExecCmd {
	editor := repo.Editor
	editCmd := NewExecCmd(editor)
	r := regexp.MustCompile("^[mg]?vim$")
	if r.MatchString(editor) {
		editCmd.WithArg("-c")
		editCmd.WithArg("set ft=gitcommit")
	}
	editCmd.WithArg(messageFile)

	return editCmd
}

func readTitleAndBodyFromFile(messageFile string) (title, body string, err error) {
	f, err := os.Open(messageFile)
	defer f.Close()
	if err != nil {
		return "", "", err
	}

	reader := bufio.NewReader(f)

	return readTitleAndBody(reader)
}

func readTitleAndBody(reader *bufio.Reader) (title, body string, err error) {
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
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}

	return string(ln), err
}
