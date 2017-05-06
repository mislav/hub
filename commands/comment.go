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
  cmdComment = &Command{
    Run: listComments,
    Usage: `
comment --issue=<issue> [-@ <USER>] [-f <FORMAT>] [-d <DATE>] [-o <SORT_KEY> [-^]]
comment create --issue=<issue> [-oc] [-m <MESSAGE>|-F <FILE>]
`,
    Long: `Manage GitHub comments for the current project.

## Commands:

With no arguments, show a list of comments associated with an issue.

  * _create_:
    Add a comment to the specified issue number.

## Options:
  -@, --mentioned <USER>
    Display only comments mentioning <USER>.

  -f, --format <FORMAT>
    Pretty print the contents of the comments using format <FORMAT> (default:
    "%sC%>(8)%i%Creset  %t%  l%n"). See the "PRETTY FORMATS" section of the
    git-log manual for some additional details on how placeholders are used in
    format. The available placeholders for comments are:

    %I: comment id

    %i: comment id prefixed with "#"

    %U: the URL of this comment

    %b: body

    %au: login name of author

    %cD: created date-only (no time of day)

    %cr: created date, relative

    %ct: created date, UNIX timestamp

    %cI: created date, ISO 8601 format

    %uD: updated date-only (no time of day)

    %ur: updated date, relative

    %ut: updated date, UNIX timestamp

    %uI: updated date, ISO 8601 format

  -m, --message <MESSAGE>
    The comment description.

  -F, --file <FILE>
    Read the comment description from <FILE>.

  -e, --edit
    Further edit the contents of <FILE> in a text editor before submitting.

  -o, --browse
    Open the new comment in a web browser.

  -c, --copy
    Put the URL of the new comment to clipboard instead of printing it.

  -d, --since <DATE>
    Display only comments updated on or after <DATE> in ISO 8601 format.

  -o, --sort <SORT_KEY>
    Sort displayed comments by "created" (default), "updated" or "comments".

  -^ --sort-ascending
    Sort by ascending dates instead of descending.
`,
  }

  cmdCreateComment = &Command{
    Key:   "create",
    Run:   createComment,
    Usage: "comment create [issue] [-o] [-m <MESSAGE>|-F <FILE>]",
    Long:  "Add a comment to the specified issue.",
  }

  flagIssueNumber,
  flagCommentFormat,
  flagCommentMessage,
  flagCommentCreator,
  flagCommentMentioned,
  flagCommentSince,
  flagCommentFile string

  flagCommentEdit,
  flagCommentCopy,
  flagCommentBrowse bool
)

func init() {
  cmdCreateComment.Flag.StringVarP(&flagIssueNumber, "issue", "i", "", "ISSUE")
  cmdCreateComment.Flag.StringVarP(&flagCommentMessage, "message", "m", "", "MESSAGE")
  cmdCreateComment.Flag.StringVarP(&flagCommentFile, "file", "F", "", "FILE")
  cmdCreateComment.Flag.BoolVarP(&flagCommentBrowse, "browse", "o", false, "BROWSE")
  cmdCreateComment.Flag.BoolVarP(&flagCommentCopy, "copy", "c", false, "COPY")
  cmdCreateComment.Flag.BoolVarP(&flagCommentEdit, "edit", "e", false, "EDIT")

  cmdComment.Flag.StringVarP(&flagIssueNumber, "issue", "i", "", "ISSUE")
  cmdComment.Flag.StringVarP(&flagCommentFormat, "format", "f", "%>(12)%i %b", "FORMAT")
  cmdComment.Flag.StringVarP(&flagCommentCreator, "creator", "c", "", "CREATOR")
  cmdComment.Flag.StringVarP(&flagCommentMentioned, "mentioned", "@", "", "USER")
  cmdComment.Flag.StringVarP(&flagCommentSince, "since", "d", "", "DATE")

  cmdComment.Use(cmdCreateComment)
  CmdRunner.Use(cmdComment)
}

func getIssueNumber(cmd *Command) int {
  if !cmd.FlagPassed("issue") {
    utils.Check(fmt.Errorf("Aborting because issue number was not provided"))
  }
  issueNumber, err := strconv.Atoi(flagIssueNumber)
  if err != nil {
    utils.Check(fmt.Errorf("Aborting because issue number was not an integer"))
  }
  return issueNumber
}

func listComments(cmd *Command, args *Args) {
  issueNumber := getIssueNumber(cmd)
  localRepo, err := github.LocalRepo()
  utils.Check(err)

  project, err := localRepo.MainProject()
  utils.Check(err)

  gh := github.NewClient(project.Host)

  if args.Noop {
    ui.Printf("Would request list of comments for %s\n", project)
  } else {
    flagFilters := map[string]string{
      "creator":   flagCommentCreator,
      "mentioned": flagCommentMentioned,
    }
    filters := map[string]interface{}{}
    for flag, filter := range flagFilters {
      if cmd.FlagPassed(flag) {
        filters[flag] = filter
      }
    }

    if cmd.FlagPassed("since") {
      if sinceTime, err := time.ParseInLocation("2006-01-02", flagCommentSince, time.Local); err == nil {
        filters["since"] = sinceTime.Format(time.RFC3339)
      } else {
        filters["since"] = flagCommentSince
      }
    }

    comments, err := gh.FetchIssueComments(project, issueNumber, filters)
    utils.Check(err)

    maxNumWidth := 0
    for _, comment := range comments {
      if numWidth := len(strconv.Itoa(comment.Id)); numWidth > maxNumWidth {
        maxNumWidth = numWidth
      }
    }

    colorize := ui.IsTerminal(os.Stdout)
    for _, comment := range comments {
      ui.Printf(formatComment(comment, flagCommentFormat, colorize))
    }
  }

  args.NoForward()
}

func formatComment(comment github.Comment, format string, colorize bool) string {
  var createdDate, createdAtISO8601, createdAtUnix, createdAtRelative,
    updatedDate, updatedAtISO8601, updatedAtUnix, updatedAtRelative string
  if !comment.CreatedAt.IsZero() {
    createdDate = comment.CreatedAt.Format("02 Jan 2006")
    createdAtISO8601 = comment.CreatedAt.Format(time.RFC3339)
    createdAtUnix = fmt.Sprintf("%d", comment.CreatedAt.Unix())
    createdAtRelative = utils.TimeAgo(comment.CreatedAt)
  }
  if !comment.UpdatedAt.IsZero() {
    updatedDate = comment.UpdatedAt.Format("02 Jan 2006")
    updatedAtISO8601 = comment.UpdatedAt.Format(time.RFC3339)
    updatedAtUnix = fmt.Sprintf("%d", comment.UpdatedAt.Unix())
    updatedAtRelative = utils.TimeAgo(comment.UpdatedAt)
  }

  placeholders := map[string]string{
    "I":  fmt.Sprintf("%d", comment.Id),
    "i":  fmt.Sprintf("#%d", comment.Id),
    "U":  comment.HtmlUrl,
    "b":  comment.Body,
    "au": comment.User.Login,
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

func createComment(cmd *Command, args *Args) {
  issueNumber := getIssueNumber(cmd)
  localRepo, err := github.LocalRepo()
  utils.Check(err)

  project, err := localRepo.MainProject()
  utils.Check(err)

  gh := github.NewClient(project.Host)

  var title string
  var body string
  var editor *github.Editor

  if cmd.FlagPassed("message") {
    title, body = readMsg(flagCommentMessage)
  } else if cmd.FlagPassed("file") {
    title, body, editor, err = readMsgFromFile(flagCommentFile, flagCommentEdit, "COMMENT", "comment")
    utils.Check(err)
  } else {
    cs := git.CommentChar()
    message := strings.Replace(fmt.Sprintf(`
# Creating a comment for %s
#
# Write a message for this comment.
`, project), "#", cs, -1)

    editor, err := github.NewEditor("COMMENT", "comment", message)
    utils.Check(err)

    title, body, err = editor.EditTitleAndBody()
    utils.Check(err)
  }

  if editor != nil {
    defer editor.DeleteFile()
  }

  if title == "" {
    utils.Check(fmt.Errorf("Aborting creation due to empty comment title"))
  }

  // TODO: Don't split the title from the body just to put it back :P
  params := map[string]interface{}{
    "body":  fmt.Sprintf("%s\n%s", title, body),
  }

  args.NoForward()
  if args.Noop {
    ui.Printf("Would create comment `%s' for %s\n", params["body"], project)
  } else {
    comment, err := gh.CreateIssueComment(project, issueNumber, params)
    utils.Check(err)

    printBrowseOrCopy(args, comment.HtmlUrl, flagCommentBrowse, flagCommentCopy)
  }
}
