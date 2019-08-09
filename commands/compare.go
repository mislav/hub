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
	Run: compare,
	Usage: `
compare [-uc] [<USER>] [[<START>...]<END>]
compare [-uc] [-b <BASE>]
`,
	Long: `Open a GitHub compare page in a web browser.

## Options:
	-u, --url
		Print the URL instead of opening it.

	-c, --copy
		Put the URL to clipboard instead of opening it.

	-b, --base <BASE>
		Base branch to compare against in case no explicit arguments were given.

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

func init() {
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

	flagCompareBase := args.Flag.Value("--base")

	// Without any flags, try to figure out something sensible to do
	if args.IsParamsEmpty() {
		// First, check we're on a branch
		localBranch, err := localRepo.CurrentBranch()
		utils.Check(err)

		// If that branch has an explicit upstream, follow it
		branch, err = localBranch.Upstream()
		if err == nil {
			// Look at the matching remote
			remote, err := localRepo.RemoteByName(branch.RemoteName())
			utils.Check(err)

			// And match it to a project
			project, err = remote.Project()
			utils.Check(err)
		} else {
			// Otherwise assume a simple push strategy, that we're
			// pushing to the same named branch on the default remote
			branch = localBranch
		}

		r = branch.ShortName()

		if flagCompareBase != "" {
			if r == flagCompareBase {
				utils.Check(command.UsageError(""))
			} else {
				r = parseCompareRange(flagCompareBase + "..." + r)
			}
		}
	} else {
		if flagCompareBase != "" {
			utils.Check(command.UsageError(""))
		} else {
			r = parseCompareRange(args.RemoveParam(args.ParamsSize() - 1))
			project, err = localRepo.CurrentProject()
			if args.IsParamsEmpty() {
				utils.Check(err)
			} else {
				projectName := ""
				projectHost := ""
				if err == nil {
					projectName = project.Name
					projectHost = project.Host
				}
				project = github.NewProject(args.RemoveParam(args.ParamsSize()-1), projectName, projectHost)
				if project.Name == "" {
					utils.Check(fmt.Errorf("Error: missing project name (owner: %q)\n", project.Owner))
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
	flagCompareURLOnly := args.Flag.Bool("--url")
	flagCompareCopy := args.Flag.Bool("--copy")
	printBrowseOrCopy(args, url, !flagCompareURLOnly && !flagCompareCopy, flagCompareCopy)
}

func parseCompareRange(r string) string {
	shaOrTag := fmt.Sprintf("((?:%s:)?\\w(?:[\\w.-]*\\w)?)", OwnerRe)
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
