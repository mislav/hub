package commands

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var cmdCompare = &Command{
	Run:   compare,
	Usage: "compare [-u] [USER] [<START>...]<END>",
	Short: "Open a compare page on GitHub",
	Long: `Open a GitHub compare view page in the system's default web browser.
<START> to <END> are branch names, tag names, or commit SHA1s specifying
the range of history to compare. If a range with two dots ("a..b") is given,
it will be transformed into one with three dots. If <START> is omitted,
GitHub will compare against the base branch (the default is "master").
If <END> is omitted, GitHub compare view is opened for the current branch.
With "-u", outputs the URL rather than opening the browser.
`,
}

var (
	flagCompareURLOnly bool
)

func init() {
	cmdCompare.Flag.BoolVarP(&flagCompareURLOnly, "url-only", "u", false, "URL only")

	CmdRunner.Use(cmdCompare)
}

/*
  $ gh compare refactor
  > open https://github.com/CURRENT_REPO/compare/refactor

  $ gh compare 1.0..1.1
  > open https://github.com/CURRENT_REPO/compare/1.0...1.1

  $ gh compare -u other-user patch
  > open https://github.com/other-user/REPO/compare/patch
*/
func compare(command *Command, args *Args) {
	localRepo, err := github.LocalRepo()
	utils.Check(err)

	var (
		branch  *github.Branch
		project *github.Project
		r       string
	)

	branch, project, err = localRepo.RemoteBranchAndProject("", false)
	utils.Check(err)

	if args.IsParamsEmpty() {
		if branch != nil && !branch.IsMaster() {
			r = branch.ShortName()
		} else {
			err = fmt.Errorf("Usage: hub compare [USER] [<START>...]<END>")
			utils.Check(err)
		}
	} else {
		r = parseCompareRange(args.RemoveParam(args.ParamsSize() - 1))
		if args.IsParamsEmpty() {
			project, err = localRepo.CurrentProject()
			utils.Check(err)
		} else {
			project = github.NewProject(args.RemoveParam(args.ParamsSize()-1), "", "")
		}
	}

	if project == nil {
		project, err = localRepo.CurrentProject()
		utils.Check(err)
	}

	subpage := utils.ConcatPaths("compare", rangeQueryEscape(r))
	url := project.WebURL("", "", subpage)
	launcher, err := utils.BrowserLauncher()
	utils.Check(err)

	if flagCompareURLOnly {
		args.Replace("echo", url)
	} else {
		args.Replace(launcher[0], "", launcher[1:]...)
		args.AppendParams(url)
	}
}

func parseCompareRange(r string) string {
	shaOrTag := fmt.Sprintf("((?:%s:)?\\w[\\w.-]+\\w)", OwnerRe)
	shaOrTagRange := fmt.Sprintf("^%s\\.\\.%s$", shaOrTag, shaOrTag)
	shaOrTagRangeRegexp := regexp.MustCompile(shaOrTagRange)
	return shaOrTagRangeRegexp.ReplaceAllString(r, "$1...$2")
}

// characters we want to allow unencoded in compare views
var compareUnescaper = strings.NewReplacer(
	"%2F", "/",
	"%3A", ":",
	"%5E", "^",
	"%7E", "~",
	"%2A", "*",
	"%21", "!",
)

func rangeQueryEscape(r string) string {
	if strings.Contains(r, "..") {
		return r
	} else {
		return compareUnescaper.Replace(url.QueryEscape(r))
	}
}
