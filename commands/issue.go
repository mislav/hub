package commands

import (
	"fmt"

	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var (
	cmdIssue = &Command{
		Run:   issue,
		Usage: "issue",
		Long:  `List open issues for the current GitHub project.`,
	}

	cmdCreateIssue = &Command{
		Key:   "create",
		Run:   createIssue,
		Usage: "issue create [-m <MESSAGE>|-f <FILE>] [-l <LABELS>]",
		Long: `File an issue for the current GitHub project.

## Options:
	-m, --message <MESSAGE>
		Use the first line of <MESSAGE> as issue title, and the rest as issue description.

	-f, --file <FILE>
		Read the issue title and description from <FILE>.

	-l, --labels <LABELS>
		Add a comma-separated list of labels to this issue.
`,
	}

	flagIssueAssignee,
	flagIssueMessage,
	flagIssueFile string

	flagIssueLabels listFlag
)

func init() {
	cmdCreateIssue.Flag.StringVarP(&flagIssueMessage, "message", "m", "", "MESSAGE")
	cmdCreateIssue.Flag.StringVarP(&flagIssueFile, "file", "f", "", "FILE")
	cmdCreateIssue.Flag.VarP(&flagIssueLabels, "label", "l", "LABEL")

	cmdIssue.Flag.StringVarP(&flagIssueAssignee, "assignee", "a", "", "ASSIGNEE")

	cmdIssue.Use(cmdCreateIssue)
	CmdRunner.Use(cmdIssue)
}

/*
  $ hub issue
*/
func issue(cmd *Command, args *Args) {
	runInLocalRepo(func(localRepo *github.GitHubRepo, project *github.Project, gh *github.Client) {
		if args.Noop {
			ui.Printf("Would request list of issues for %s\n", project)
		} else {
			issues, err := gh.Issues(project)
			utils.Check(err)
			for _, issue := range issues {
				var url string
				// use the pull request URL if we have one
				if issue.PullRequest.HTMLURL != "" {
					url = issue.PullRequest.HTMLURL
				} else {
					url = issue.HTMLURL
				}

				if flagIssueAssignee == "" || issue.Assignee.Login == flagIssueAssignee {
					// "nobody" should have more than 1 million github issues
					ui.Printf("% 7d] %s ( %s )\n", issue.Number, issue.Title, url)
				}
			}
		}
	})
}

func createIssue(cmd *Command, args *Args) {
	runInLocalRepo(func(localRepo *github.GitHubRepo, project *github.Project, gh *github.Client) {
		if args.Noop {
			ui.Printf("Would create an issue for %s\n", project)
		} else {
			title, body, err := getTitleAndBodyFromFlags(flagIssueMessage, flagIssueFile)
			utils.Check(err)

			if title == "" {
				title, body, err = writeIssueTitleAndBody(project)
				utils.Check(err)
			}

			issue, err := gh.CreateIssue(project, title, body, flagIssueLabels)
			utils.Check(err)

			ui.Println(issue.HTMLURL)
		}
	})
}

func writeIssueTitleAndBody(project *github.Project) (string, string, error) {
	message := `
# Creating issue for %s.
#
# Write a message for this issue. The first block of
# text is the title and the rest is the description.
`
	message = fmt.Sprintf(message, project.Name)

	editor, err := github.NewEditor("ISSUE", "issue", message)
	if err != nil {
		return "", "", err
	}

	defer editor.DeleteFile()

	return editor.EditTitleAndBody()
}
