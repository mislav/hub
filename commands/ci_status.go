package commands

import (
	"fmt"
	"os"
	"sort"

	"github.com/github/hub/v2/git"
	"github.com/github/hub/v2/github"
	"github.com/github/hub/v2/ui"
	"github.com/github/hub/v2/utils"
)

var cmdCiStatus = &Command{
	Run:   ciStatus,
	Usage: "ci-status [-v] [<COMMIT>]",
	Long: `Display status of GitHub checks for a commit.

## Options:
	-v, --verbose
		Print detailed report of all status checks and their URLs.

	-f, --format <FORMAT>
		Pretty print all status checks using <FORMAT> (implies ''--verbose''). See the
		"PRETTY FORMATS" section of git-log(1) for some additional details on how
		placeholders are used in format. The available placeholders for checks are:

		%U: the URL of this status check

		%S: check state (e.g. "success", "failure")

		%sC: set color to red, green, or yellow, depending on state

		%t: name of the status check

	--color[=<WHEN>]
		Enable colored output even if stdout is not a terminal. <WHEN> can be one
		of "always" (default for ''--color''), "never", or "auto" (default).

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

var severityList []string

func init() {
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

		state := ""
		if len(response.Statuses) > 0 {
			for _, status := range response.Statuses {
				if checkSeverity(status.State) > checkSeverity(state) {
					state = status.State
				}
			}
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

		verbose := args.Flag.Bool("--verbose") || args.Flag.HasReceived("--format")
		if verbose && len(response.Statuses) > 0 {
			colorize := colorizeOutput(args.Flag.HasReceived("--color"), args.Flag.Value("--color"))
			ciVerboseFormat(response.Statuses, args.Flag.Value("--format"), colorize)
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

func ciVerboseFormat(statuses []github.CIStatus, formatString string, colorize bool) {
	contextWidth := 0
	for _, status := range statuses {
		if len(status.Context) > contextWidth {
			contextWidth = len(status.Context)
		}
	}

	sort.SliceStable(statuses, func(a, b int) bool {
		return stateRank(statuses[a].State) < stateRank(statuses[b].State)
	})

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

		placeholders := map[string]string{
			"S":  status.State,
			"sC": "",
			"t":  status.Context,
			"U":  status.TargetURL,
		}

		if colorize {
			placeholders["sC"] = fmt.Sprintf("\033[%dm", color)
		}

		format := formatString
		if format == "" {
			if status.TargetURL == "" {
				format = fmt.Sprintf("%%sC%s%%Creset\t%%t\n", stateMarker)
			} else {
				format = fmt.Sprintf("%%sC%s%%Creset\t%%<(%d)%%t\t%%U\n", stateMarker, contextWidth)
			}
		}
		ui.Print(ui.Expand(format, placeholders, colorize))
	}
}

func stateRank(state string) uint32 {
	switch state {
	case "failure", "error", "action_required", "cancelled", "timed_out":
		return 1
	case "success", "neutral":
		return 3
	default:
		return 2
	}
}
