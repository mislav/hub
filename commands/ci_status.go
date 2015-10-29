package commands

import (
	"fmt"
	"os"

	"github.com/github/hub/git"
	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var cmdCiStatus = &Command{
	Run:   ciStatus,
	Usage: "ci-status [-v] [COMMIT]",
	Short: "Show CI status of a commit",
	Long: `Looks up the SHA for <COMMIT> in GitHub Status API and displays the latest
status. Exits with one of:
success (0), error (1), failure (1), pending (2), no status (3)

If "-v" is given, additionally print detailed report of all checks and their URLs.
`,
}

var flagCiStatusVerbose bool

func init() {
	cmdCiStatus.Flag.BoolVarP(&flagCiStatusVerbose, "verbose", "v", false, "VERBOSE")

	CmdRunner.Use(cmdCiStatus)
}

/*
  $ gh ci-status
  > (prints CI state of HEAD and exits with appropriate code)

  $ gh ci-status -v
  > (prints CI states and URLs to CI build results for HEAD and exits with appropriate code)

  $ gh ci-status BRANCH
  > (prints CI state of BRANCH and exits with appropriate code)

  $ gh ci-status SHA
  > (prints CI state of SHA and exits with appropriate code)
*/
func ciStatus(cmd *Command, args *Args) {
	ref := "HEAD"
	if !args.IsParamsEmpty() {
		ref = args.RemoveParam(0)
	}

	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.MainProject()
	utils.Check(err)

	sha, err := git.Ref(ref)
	if err != nil {
		err = fmt.Errorf("Aborted: no revision could be determined from '%s'", ref)
	}
	utils.Check(err)

	if args.Noop {
		ui.Printf("Would request CI status for %s\n", sha)
	} else {
		gh := github.NewClient(project.Host)
		response, err := gh.FetchCIStatus(project, sha)
		utils.Check(err)

		state := response.State
		if len(response.Statuses) == 0 {
			state = ""
		}

		var exitCode int
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

		targetUrl := relevantTargetUrl(response.Statuses, state)

		if flagCiStatusVerbose && targetUrl != "" {
			ui.Printf("%s: %s\n", state, targetUrl)
		} else {
			if state != "" {
				ui.Println(state)
			} else {
				ui.Println("no status")
			}
		}

		os.Exit(exitCode)
	}
}

func relevantTargetUrl(statuses []github.CIStatus, state string) string {
	for _, status := range statuses {
		if status.State == state {
			return status.TargetUrl
		}
	}
	return ""
}

func verboseFormat(statuses []github.CIStatus) {
	colorize := github.IsTerminal(os.Stdout)

	contextWidth := 0
	for _, status := range statuses {
		if len(status.Context) > contextWidth {
			contextWidth = len(status.Context)
		}
	}

	for _, status := range statuses {
		var color int
		var stateMarker string
		switch status.State {
		case "success":
			stateMarker = "✔︎"
			color = 32
		case "failure", "error":
			stateMarker = "✖︎"
			color = 31
		case "pending":
			stateMarker = "●"
			color = 33
		}

		if colorize {
			stateMarker = fmt.Sprintf("\033[%dm%s\033[0m", color, stateMarker)
		}

		if status.TargetUrl == "" {
			ui.Printf("%s\t%s\n", stateMarker, status.Context)
		} else {
			ui.Printf("%s\t%-*s\t%s\n", stateMarker, contextWidth, status.Context, status.TargetUrl)
		}
	}
}
