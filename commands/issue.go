package commands

import (
	"fmt"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"os"
)

var cmdIssue = &Command{
	Run:   issue,
	Usage: "issue",
	Short: "Manipulate issues on GitHub",
	Long:  `Lists summary of the open issues for the project that the "origin" remove points to.`,
}

/*
  $ gh issue
*/
func issue(cmd *Command, args *Args) {

	gh := github.New()

	// list all issues

	if args.Noop {
		fmt.Printf("Would request list of issues for %s\n", gh.Project)
	} else {
		issues, err := gh.Issues()
		utils.Check(err)
		for _, issue := range issues {
			var url string
			// use the pull request URL if we have one
			if issue.PullRequest.HTMLURL != "" {
				url = issue.PullRequest.HTMLURL
			} else {
				url = issue.HTMLURL
			}
			// "nobody" should have more than 1 million github issues
			fmt.Printf("% 7d] %s ( %s )\n", issue.Number, issue.Title, url)
		}
	}

	os.Exit(0)
}
