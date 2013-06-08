package commands

import (
	"fmt"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"os"
)

var cmdCi = &Command{
	Run:   ci,
	Usage: "ci [COMMIT]",
	Short: "Show CI status of a commit",
	Long: `Looks up the SHA for COMMIT in GitHub Status API and displays the latest
status. If no COMMIT is provided, HEAD will be used. Exits with one of:

success (0), error (1), failure (1), pending (2), no status (3)
`,
}

func ci(cmd *Command, args []string) {
	ref := "HEAD"
	if len(args) > 0 {
		ref = args[0]
	}

	ref, err := git.Ref(ref)
	utils.Check(err)

	gh := github.New()
	status, err := gh.CiStatus(ref)
	utils.Check(err)

	var state string
	var targetUrl string
	var desc string
	var exitCode int
	if status == nil {
		state = "no status"
	} else {
		state = status.State
		targetUrl = status.TargetUrl
		desc = status.Description
	}

	switch state {
	case "success":
		exitCode = 0
	case "failure", "error":
		exitCode = 1
	case "pending":
		exitCode = 2
	default:
		exitCode = 3
	}

	fmt.Println(state)
	if targetUrl != "" {
		fmt.Println(targetUrl)
	}
	if desc != "" {
		fmt.Println(desc)
	}

	os.Exit(exitCode)
}
