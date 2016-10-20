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
	Usage: "compare [-u] [<USER>] [[<START>...]<END>]",
	Long: `Open a GitHub compare page in a web browser.

## Options:
	-u
		Print the URL instead of opening it.

	[<START>...]<END>
		Branch names, tag names, or commit SHAs specifying the range to compare.
		<END> defaults to the current branch name.

		If a range with two dots ('A..B') is given, it will be transformed into a
		range with three dots.

## Examples:
		$ hub compare refactor
		> open https://github.com/USER/REPO/compare/refactor

		$ hub compare v1.0..v1.1
		> open https://github.com/USER/REPO/compare/v1.0...v1.1

		$ hub compare -u jingweno feature
		> echo https://github.com/jingweno/REPO/compare/feature

## See also:

hub-browse(1), hub(1)
`,
}

var (
	flagCompareURLOnly bool
	flagCompareBase    string
)

func init() {
	cmdCompare.Flag.BoolVarP(&flagCompareURLOnly, "url-only", "u", false, "URL only")
	cmdCompare.Flag.StringVarP(&flagCompareBase, "base", "b", "", "BASE")

	CmdRunner.Use(cmdCompare)
}

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

	usageHelp := func() {
		utils.Check(fmt.Errorf("Usage: hub compare [-u] [-b <BASE>] [<USER>] [[<START>...]<END>]"))
	}

	if args.IsParamsEmpty() {
		if branch == nil ||
			(branch.IsMaster() && flagCompareBase == "") ||
			(flagCompareBase == branch.ShortName()) {

			usageHelp()
		} else {
			r = branch.ShortName()
			if flagCompareBase != "" {
				r = parseCompareRange(flagCompareBase + "..." + r)
			}
		}
	} else {
		if flagCompareBase != "" {
			usageHelp()
		} else {
			r = parseCompareRange(args.RemoveParam(args.ParamsSize() - 1))
			project, err = localRepo.CurrentProject()
			if args.IsParamsEmpty() {
				utils.Check(err)
			} else {
				projectName := ""
				if err == nil {
					projectName = project.Name
				}
				project = github.NewProject(args.RemoveParam(args.ParamsSize()-1), projectName, "")
				if project.Name == "" {
					utils.Check(fmt.Errorf("error: missing project name (owner: %q)\n", project.Owner))
				}
			}
		}
	}

	if project == nil {
		project, err = localRepo.CurrentProject()
		utils.Check(err)
	}

	subpage := utils.ConcatPaths("compare", rangeQueryEscape(r))
	url := project.WebURL("", "", subpage)

	args.NoForward()
	printBrowseOrCopy(args, url, !flagCompareURLOnly, false)
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
