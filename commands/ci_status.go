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
	Usage: "ci-status [-v] [<COMMIT>]",
	Long: `Display status of GitHub checks for a commit.

## Options:
	-v
		Print detailed report of all status checks and their URLs.

	<COMMIT>
		A commit SHA or branch name (default: "HEAD").

Possible outputs and exit statuses:

- success, neutral: 0
- failure, error, action_required, cancelled, timed_out: 1
- pending: 2

## See also:

hub-pull-request(1), hub(1)
`,
}

var flagCiStatusVerbose bool
var severityList []string

func init() {
	cmdCiStatus.Flag.BoolVarP(&flagCiStatusVerbose, "verbose", "v", false, "VERBOSE")

	CmdRunner.Use(cmdCiStatus)

	severityList = []string{
		"neutral",
		"success",
		"pending",
		"cancelled",
		"timed_out",
		"action_required",
		"failure",
		"error",
	}
}

func checkSeverity(targetState string) int {
	for i, state := range severityList {
		if state == targetState {
			return i
		}
	}
	return -1
}

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
		if len(response.Statuses) > 0 {
			for _, status := range response.Statuses {
				if checkSeverity(status.State) > checkSeverity(state) {
					state = status.State
				}
			}
		} else if len(response.Statuses) == 0 {
			state = ""
		}

		var exitCode int
		switch state {
		case "success", "neutral":
			exitCode = 0
		case "failure", "error", "action_required", "cancelled", "timed_out":
			exitCode = 1
		case "pending":
			exitCode = 2
		default:
			exitCode = 3
		}

		if flagCiStatusVerbose && len(response.Statuses) > 0 {
			verboseFormat(response.Statuses)
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

func verboseFormat(statuses []github.CIStatus) {
	colorize := ui.IsTerminal(os.Stdout)

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
		case "failure", "error", "action_required", "cancelled", "timed_out":
			stateMarker = "✖︎"
			color = 31
		case "neutral":
			stateMarker = "◦"
			color = 30
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
