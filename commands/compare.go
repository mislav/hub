package commands

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/github/hub/v2/github"
	"github.com/github/hub/v2/utils"
)

var cmdCompare = &Command{
	Run: compare,
	Usage: `
compare [-uc] [-b <BASE>]
compare [-uc] [<OWNER>] [<BASE>...]<HEAD>
`,
	Long: `Open a GitHub compare page in a web browser.

## Options:
	-u, --url
		Print the URL instead of opening it.

	-c, --copy
		Put the URL to clipboard instead of opening it.

	-b, --base <BASE>
		Base branch to compare against in case no explicit arguments were given.

	[<BASE>...]<HEAD>
		Branch names, tag names, or commit SHAs specifying the range to compare.
		If a range with two dots (''A..B'') is given, it will be transformed into a
		range with three dots.

		The <BASE> portion defaults to the default branch of the repository.

		The <HEAD> argument defaults to the current branch. If the current branch
		is not pushed to a remote, the command will error.

	<OWNER>
		Optionally specify the owner of the repository for the compare page URL.

## Examples:
		$ hub compare
		> open https://github.com/OWNER/REPO/compare/BRANCH

		$ hub compare refactor
		> open https://github.com/OWNER/REPO/compare/refactor

		$ hub compare v1.0..v1.1
		> open https://github.com/OWNER/REPO/compare/v1.0...v1.1

		$ hub compare -u jingweno feature
		https://github.com/jingweno/REPO/compare/feature

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

	mainProject, err := localRepo.MainProject()
	utils.Check(err)

	host, err := github.CurrentConfig().PromptForHost(mainProject.Host)
	utils.Check(err)

	var r string
	flagCompareBase := args.Flag.Value("--base")

	if args.IsParamsEmpty() {
		currentBranch, err := localRepo.CurrentBranch()
		if err != nil {
			utils.Check(command.UsageError(err.Error()))
		}

		var remoteBranch *github.Branch
		var remoteProject *github.Project

		remoteBranch, remoteProject, err = findPushTarget(currentBranch)
		if err != nil {
			if remoteProject, err = deducePushTarget(currentBranch, host.User); err == nil {
				remoteBranch = currentBranch
			} else {
				utils.Check(fmt.Errorf("the current branch '%s' doesn't seem pushed to a remote", currentBranch.ShortName()))
			}
		}

		r = remoteBranch.ShortName()
		if remoteProject.SameAs(mainProject) {
			if flagCompareBase == "" && remoteBranch.IsMaster() {
				utils.Check(fmt.Errorf("the branch to compare '%s' is the default branch", remoteBranch.ShortName()))
			}
		} else {
			r = fmt.Sprintf("%s:%s", remoteProject.Owner, r)
		}

		if flagCompareBase == r {
			utils.Check(fmt.Errorf("the branch to compare '%s' is the same as --base", r))
		} else if flagCompareBase != "" {
			r = fmt.Sprintf("%s...%s", flagCompareBase, r)
		}
	} else {
		if flagCompareBase != "" {
			utils.Check(command.UsageError(""))
		} else {
			r = parseCompareRange(args.RemoveParam(args.ParamsSize() - 1))
			if !args.IsParamsEmpty() {
				owner := args.RemoveParam(args.ParamsSize() - 1)
				mainProject = github.NewProject(owner, mainProject.Name, mainProject.Host)
			}
		}
	}

	url := mainProject.WebURL("", "", "compare/"+rangeQueryEscape(r))

	args.NoForward()
	flagCompareURLOnly := args.Flag.Bool("--url")
	flagCompareCopy := args.Flag.Bool("--copy")
	printBrowseOrCopy(args, url, !flagCompareURLOnly && !flagCompareCopy, flagCompareCopy)
}

func parseCompareRange(r string) string {
	shaOrTag := fmt.Sprintf("((?:%s:)?\\w(?:[\\w/.-]*\\w)?)", OwnerRe)
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
	}
	return compareUnescaper.Replace(url.QueryEscape(r))
}
