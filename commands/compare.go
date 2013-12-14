package commands

import (
	"fmt"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"regexp"
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
	cmdCompare.Flag.BoolVar(&flagCompareURLOnly, "u", false, "URL only")
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
	localRepo := github.LocalRepo()
	var (
		project *github.Project
		r       string
		err     error
	)
	if args.IsParamsEmpty() {
		current, err := localRepo.CurrentBranch()
		utils.Check(err)

		var printUsage bool
		upstream, err := current.Upstream()

		if err == nil {
			master, err := localRepo.MasterBranch()
			if err != nil || master.ShortName() == upstream.ShortName() {
				printUsage = true
			}
		} else {
			printUsage = true
		}

		if printUsage {
			err = fmt.Errorf(command.FormattedUsage())
			utils.Check(err)
		} else {
			r = upstream.ShortName()
			project, err = localRepo.CurrentProject()
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

	subpage := utils.ConcatPaths("compare", r)
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
