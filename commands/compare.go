package commands

import (
	"fmt"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"regexp"
)

var cmdCompare = &Command{
	Run:   compare,
	Usage: "compare [-u USER] [START...] END",
	Short: "Open a compare page on GitHub",
	Long: `Open a GitHub compare view page in the system's default web browser.
START to END are branch names, tag names, or commit SHA1s specifying
the range of history to compare. If a range with two dots (a..b) is given,
it will be transformed into one with three dots. If START is omitted,
GitHub will compare against the base branch (the default is "master").
`,
}

var flagCompareUser string

func init() {
	cmdCompare.Flag.StringVar(&flagCompareUser, "u", "", "USER")
}

func compare(command *Command, args []string) {
	project := github.CurrentProject()

	var r string
	if len(args) == 0 {
		repo := project.LocalRepo()
		r = repo.Head
	} else {
		r = args[0]
	}

	r = transformToTripleDots(r)
	subpage := utils.ConcatPaths("compare", r)
	url := project.WebURL("", flagCompareUser, subpage)
	err := browserCommand(url)
	utils.Check(err)
}

func transformToTripleDots(r string) string {
	ownerRe := "[a-zA-Z0-9][a-zA-Z0-9-]*"
	shaOrTag := fmt.Sprintf("((?:%s:)?\\w[\\w.-]+\\w)", ownerRe)
	shaOrTagRange := fmt.Sprintf("^%s\\.\\.%s$", shaOrTag, shaOrTag)
	shaOrTagRangeRegexp := regexp.MustCompile(shaOrTagRange)
	return shaOrTagRangeRegexp.ReplaceAllString(r, "$1...$2")
}
