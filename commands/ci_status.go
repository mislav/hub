package commands

import (
	"fmt"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"os"
)

var cmdCiStatus = &Command{
	Run:   ciStatus,
	Usage: "ci-status [COMMIT]",
	Short: "Show CI status of a commit",
	Long: `Looks up the SHA for COMMIT in GitHub Status API and displays the latest
status. If no COMMIT is provided, HEAD will be used. Exits with one of:

success (0), error (1), failure (1), pending (2), no status (3)
`,
}

/*
  $ gh ci-status
  > (prints CI state of HEAD and exits with appropriate code)
  > One of: success (0), error (1), failure (1), pending (2), no status (3)

  $ gh ci-status BRANCH
  > (prints CI state of BRANCH and exits with appropriate code)
  > One of: success (0), error (1), failure (1), pending (2), no status (3)

  $ gh ci-status SHA
  > (prints CI state of SHA and exits with appropriate code)
  > One of: success (0), error (1), failure (1), pending (2), no status (3)
*/
func ciStatus(cmd *Command, args *Args) {
	ref := "HEAD"
	if !args.IsParamsEmpty() {
		ref = args.RemoveParam(0)
	}

	ref, err := git.Ref(ref)
	utils.Check(err)

	args.Replace("", "")
	if args.Noop {
		fmt.Printf("Would request CI status for %s", ref)
	} else {
		state, targetURL, desc, exitCode, err := fetchCiStatus(ref)
		utils.Check(err)
		fmt.Println(state)
		if targetURL != "" {
			fmt.Println(targetURL)
		}
		if desc != "" {
			fmt.Println(desc)
		}

		os.Exit(exitCode)
	}
}

func fetchCiStatus(ref string) (state, targetURL, desc string, exitCode int, err error) {
	gh := github.New()
	status, err := gh.CiStatus(ref)
	if err != nil {
		return
	}

	if status == nil {
		state = "no status"
	} else {
		state = status.State
		targetURL = status.TargetURL
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

	return
}
