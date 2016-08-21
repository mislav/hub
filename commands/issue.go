package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

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
issue create [-o] [-m <MESSAGE>|-F <FILE>] [-a <USERS>] [-M <MILESTONE>] [-l <LABELS>]
`,
		Long: `Manage GitHub issues for the current project.

## Commands:

With no arguments, show a list of open issues.

	* _create_:
		Open an issue in the current project.

## Options:
	-a, --assignee <ASSIGNEE>
		Display only issues assigned to <ASSIGNEE>.

		When opening an issue, this can be a comma-separated list of people to
		assign to the new issue.

	-s, --state <STATE>
		Display issues with state <STATE> (default: "open").

	-f, --format <FORMAT>
		Pretty print the contents of the issues using format <FORMAT> (default:
		"%sC%>(8)%i%Creset  %t%  l%n"). See the "PRETTY FORMATS" section of the
		git-log manual for some additional details on how placeholders are used in
		format. The available placeholders for issues are:

			· %I: issue number

			· %i: issue number prefixed with "#"

			· %U: the URL of this issue

			· %S: state (i.e. "open", "closed")

			· %sC: set color to red or green, depending on issue state.

			· %t: title

			· %l: colored labels

			· %L: raw, comma-separated labels

			· %b: body

			· %au: login name of author

			· %as: comma-separated list of assignees

			· %Mn: milestone number

			· %Mt: milestone title

			· %NC: number of comments

			· %Nc: number of comments wrapped in parentheses, or blank string if zero.

			· %cD: created date-only (no time of day)

			· %cr: created date, relative

			· %ct: created date, UNIX timestamp

			· %cI: created date, ISO 8601 format

			· %uD: updated date-only (no time of day)

			· %ur: updated date, relative

			· %ut: updated date, UNIX timestamp

			· %uI: updated date, ISO 8601 format

	-m, --message <MESSAGE>
		Use the first line of <MESSAGE> as issue title, and the rest as issue description.

	-F, --file <FILE>
		Read the issue title and description from <FILE>.

	-o, --browse
		Open the new issue in a web browser.

	-M, --milestone <ID>
		Add this pull request to a GitHub milestone with id <ID>.

	-l, --labels <LABELS>
		Add a comma-separated list of labels to this issue.
`,
	}

	cmdCreateIssue = &Command{
		Key:   "create",
		Run:   createIssue,
		Usage: "issue create [-o] [-m <MESSAGE>|-F <FILE>] [-a <USERS>] [-M <MILESTONE>] [-l <LABELS>]",
		Long:  "Open an issue in the current project.",
	}

	flagIssueAssignee,
	flagIssueState,
	flagIssueFormat,
	flagIssueMessage,
	flagIssueFile string

	flagIssueBrowse bool

	flagIssueMilestone uint64

	flagIssueAssignees,
	flagIssueLabels listFlag
)

func init() {
	cmdCreateIssue.Flag.StringVarP(&flagIssueMessage, "message", "m", "", "MESSAGE")
	cmdCreateIssue.Flag.StringVarP(&flagIssueFile, "file", "F", "", "FILE")
	cmdCreateIssue.Flag.Uint64VarP(&flagIssueMilestone, "milestone", "M", 0, "MILESTONE")
	cmdCreateIssue.Flag.VarP(&flagIssueLabels, "label", "l", "LABEL")
	cmdCreateIssue.Flag.VarP(&flagIssueAssignees, "assign", "a", "ASSIGNEE")
	cmdCreateIssue.Flag.BoolVarP(&flagIssueBrowse, "browse", "o", false, "BROWSE")

	cmdIssue.Flag.StringVarP(&flagIssueAssignee, "assignee", "a", "", "ASSIGNEE")
	cmdIssue.Flag.StringVarP(&flagIssueState, "state", "s", "", "STATE")
	cmdIssue.Flag.StringVarP(&flagIssueFormat, "format", "f", "%sC%>(8)%i%Creset  %t%  l%n", "FORMAT")

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
	var stateColorSwitch string
	if colorize {
		issueColor := 32
		if issue.State == "closed" {
			issueColor = 31
		}
		stateColorSwitch = fmt.Sprintf("\033[%dm", issueColor)
	}

	var labelStrings []string
	var rawLabels []string
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
		rawLabels = append(rawLabels, label.Name)
	}

	var assignees []string
	for _, assignee := range issue.Assignees {
		assignees = append(assignees, assignee.Login)
	}

	var milestoneNumber, milestoneTitle string
	if issue.Milestone != nil {
		milestoneNumber = fmt.Sprintf("%d", issue.Milestone.Number)
		milestoneTitle = issue.Milestone.Title
	}

	var numCommentsWrapped string
	numComments := fmt.Sprintf("%d", issue.Comments)
	if issue.Comments > 0 {
		numCommentsWrapped = fmt.Sprintf("(%d)", issue.Comments)
	}

	var createdDate, createdAtISO8601, createdAtUnix, createdAtRelative,
		updatedDate, updatedAtISO8601, updatedAtUnix, updatedAtRelative string
	if !issue.CreatedAt.IsZero() {
		createdDate = issue.CreatedAt.Format("02 Jan 2006")
		createdAtISO8601 = issue.CreatedAt.Format(time.RFC3339)
		createdAtUnix = fmt.Sprintf("%d", issue.CreatedAt.Unix())
		createdAtRelative = utils.TimeAgo(issue.CreatedAt)
	}
	if !issue.UpdatedAt.IsZero() {
		updatedDate = issue.UpdatedAt.Format("02 Jan 2006")
		updatedAtISO8601 = issue.UpdatedAt.Format(time.RFC3339)
		updatedAtUnix = fmt.Sprintf("%d", issue.UpdatedAt.Unix())
		updatedAtRelative = utils.TimeAgo(issue.UpdatedAt)
	}

	placeholders := map[string]string{
		"I":  fmt.Sprintf("%d", issue.Number),
		"i":  fmt.Sprintf("#%d", issue.Number),
		"U":  issue.HtmlUrl,
		"S":  issue.State,
		"sC": stateColorSwitch,
		"t":  issue.Title,
		"l":  strings.Join(labelStrings, " "),
		"L":  strings.Join(rawLabels, ", "),
		"b":  issue.Body,
		"au": issue.User.Login,
		"as": strings.Join(assignees, ", "),
		"Mn": milestoneNumber,
		"Mt": milestoneTitle,
		"NC": numComments,
		"Nc": numCommentsWrapped,
		"cD": createdDate,
		"cI": createdAtISO8601,
		"ct": createdAtUnix,
		"cr": createdAtRelative,
		"uD": updatedDate,
		"uI": updatedAtISO8601,
		"ut": updatedAtUnix,
		"ur": updatedAtRelative,
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

		if template := github.GetIssueTemplate(); template != "" {
			message = template + "\n" + message
		}

		editor, err := github.NewEditor("ISSUE", "issue", message)
		utils.Check(err)

		title, body, err = editor.EditTitleAndBody()
		utils.Check(err)
	}

	if title == "" {
		utils.Check(fmt.Errorf("Aborting creation due to empty issue title"))
	}

	params := map[string]interface{}{
		"title":     title,
		"body":      body,
		"labels":    flagIssueLabels,
		"assignees": flagIssueAssignees,
	}

	if flagIssueMilestone > 0 {
		params["milestone"] = flagIssueMilestone
	}

	if args.Noop {
		ui.Printf("Would create issue `%s' for %s\n", params["title"], project)
		os.Exit(0)
	} else {
		issue, err := gh.CreateIssue(project, params)
		utils.Check(err)

		if editor != nil {
			editor.DeleteFile()
		}

		if flagIssueBrowse {
			launcher, err := utils.BrowserLauncher()
			utils.Check(err)
			args.Replace(launcher[0], "", launcher[1:]...)
			args.AppendParams(issue.HtmlUrl)
		} else {
			ui.Println(issue.HtmlUrl)
			os.Exit(0)
		}
	}
}
