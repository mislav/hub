package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/github/hub/git"
	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var (
	cmdIssue = &Command{
		Run: listIssues,
		Usage: `
issue [-a <ASSIGNEE>] [-s <STATE>] [-f <FORMAT>]
issue create [-m <MESSAGE>|-F <FILE>] [-l <LABELS>]
`,
		Long: `Manage GitHub issues for the current project.

## Options:
	-a, --assignee <ASSIGNEE>
		Display only issues assigned to <ASSIGNEE>.

	-s, --state <STATE>
		Display issues with state <STATE> (default: "open").

	-f, --format <FORMAT>
		Pretty print the contents of the issues using format <FORMAT> (default:
		"%sC%>(8)%ih%Creset  %t%  l%n"). See the "PRETTY FORMATS" section of the
		git-log manual for some additional details on how placeholders are used in
		format. The available placeholders for issues are:

			· %in: the number of the issue.

			· %ih: the number of the issue prefixed with #.

			· %st: the state of the issue as a text (i.e. "open", "closed").

			· %sC: switch color to red if issue is closed or green if issue is open.

			· %t: the title of the issue.

			· %l: the colored labels of the issue.

			· %b: the body of the issue.

			· %u: the login of the user that opened the issue.

			· %a: the login of the user that the issue is assigned to.

	-m, --message <MESSAGE>
		Use the first line of <MESSAGE> as issue title, and the rest as issue description.

	-F, --file <FILE>
		Read the issue title and description from <FILE>.

	-l, --labels <LABELS>
		Add a comma-separated list of labels to this issue.
`,
	}

	cmdCreateIssue = &Command{
		Key:   "create",
		Run:   createIssue,
		Usage: "issue create [-m <MESSAGE>|-f <FILE>] [-l <LABELS>]",
		Long:  "File an issue for the current GitHub project.",
	}

	flagIssueAssignee,
	flagIssueState,
	flagIssueFormat,
	flagIssueMessage,
	flagIssueFile string

	flagIssueLabels listFlag
)

func init() {
	cmdCreateIssue.Flag.StringVarP(&flagIssueMessage, "message", "m", "", "MESSAGE")
	cmdCreateIssue.Flag.StringVarP(&flagIssueFile, "file", "F", "", "FILE")
	cmdCreateIssue.Flag.VarP(&flagIssueLabels, "label", "l", "LABEL")

	cmdIssue.Flag.StringVarP(&flagIssueAssignee, "assignee", "a", "", "ASSIGNEE")
	cmdIssue.Flag.StringVarP(&flagIssueState, "state", "s", "", "ASSIGNEE")
	cmdIssue.Flag.StringVarP(&flagIssueFormat, "format", "f", "%sC%>(8)%ih%Creset  %t%  l%n", "FORMAT")

	cmdIssue.Use(cmdCreateIssue)
	CmdRunner.Use(cmdIssue)
}

func listIssues(cmd *Command, args *Args) {
	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.MainProject()
	utils.Check(err)

	gh := github.NewClient(project.Host)

	if args.Noop {
		ui.Printf("Would request list of issues for %s\n", project)
	} else {
		filters := map[string]interface{}{}
		if cmd.FlagPassed("state") {
			filters["state"] = flagIssueState
		}
		if cmd.FlagPassed("assignee") {
			filters["assignee"] = flagIssueAssignee
		}

		issues, err := gh.FetchIssues(project, filters)
		utils.Check(err)

		maxNumWidth := 0
		for _, issue := range issues {
			if numWidth := len(strconv.Itoa(issue.Number)); numWidth > maxNumWidth {
				maxNumWidth = numWidth
			}
		}

		colorize := ui.IsTerminal(os.Stdout)
		for _, issue := range issues {
			if issue.PullRequest != nil {
				continue
			}

			ui.Printf(formatIssue(issue, flagIssueFormat, colorize))
		}
	}

	os.Exit(0)
}

func formatIssue(issue github.Issue, format string, colorize bool) string {
	var assigneeLogin string
	if a := issue.Assignee; a != nil {
		assigneeLogin = a.Login
	}

	var stateColorSwitch string
	if colorize {
		issueColor := 32
		if issue.State == "closed" {
			issueColor = 31
		}
		stateColorSwitch = fmt.Sprintf("\033[%dm", issueColor)
	}

	var labelStrings []string
	for _, label := range issue.Labels {
		if !colorize {
			labelStrings = append(labelStrings, fmt.Sprintf(" %s ", label.Name))
			continue
		}
		color, err := utils.NewColor(label.Color)
		if err != nil {
			utils.Check(err)
		}

		textColor := 16
		if color.Brightness() < 0.65 {
			textColor = 15
		}

		labelStrings = append(labelStrings, fmt.Sprintf("\033[38;5;%d;48;2;%d;%d;%dm %s \033[m", textColor, color.Red, color.Green, color.Blue, label.Name))
	}

	placeholders := map[string]string{
		"in": fmt.Sprintf("%d", issue.Number),
		"ih": fmt.Sprintf("#%d", issue.Number),
		"st": issue.State,
		"sC": stateColorSwitch,
		"t":  issue.Title,
		"l":  strings.Join(labelStrings, " "),
		"b":  issue.Body,
		"u":  issue.User.Login,
		"a":  assigneeLogin,
	}

	return ui.Expand(format, placeholders, colorize)
}

func createIssue(cmd *Command, args *Args) {
	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.MainProject()
	utils.Check(err)

	gh := github.NewClient(project.Host)

	var title string
	var body string
	var editor *github.Editor

	if cmd.FlagPassed("message") {
		title, body = readMsg(flagIssueMessage)
	} else if cmd.FlagPassed("file") {
		title, body, err = readMsgFromFile(flagIssueFile)
		utils.Check(err)
	} else {
		cs := git.CommentChar()
		message := strings.Replace(fmt.Sprintf(`
# Creating an issue for %s
#
# Write a message for this issue. The first block of
# text is the title and the rest is the description.
`, project), "#", cs, -1)

		editor, err := github.NewEditor("ISSUE", "issue", message)
		utils.Check(err)

		title, body, err = editor.EditTitleAndBody()
		utils.Check(err)
	}

	if title == "" {
		utils.Check(fmt.Errorf("Aborting creation due to empty issue title"))
	}

	params := &github.IssueParams{
		Title:  title,
		Body:   body,
		Labels: flagIssueLabels,
	}

	if args.Noop {
		ui.Printf("Would create issue `%s' for %s\n", params.Title, project)
	} else {
		issue, err := gh.CreateIssue(project, params)
		utils.Check(err)

		if editor != nil {
			editor.DeleteFile()
		}

		ui.Println(issue.HtmlUrl)
	}

	os.Exit(0)
}
