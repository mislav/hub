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

If "-v" is given, additionally print the URL to CI build results.
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
  > One of: success (0), error (1), failure (1), pending (2), no status (3)

  $ gh ci-status -v
  > (prints CI state of HEAD, the URL to the CI build results and exits with appropriate code)
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
		state, targetURL, exitCode, err := fetchCiStatus(project, sha)
		utils.Check(err)
		if flagCiStatusVerbose && targetURL != "" {
			ui.Printf("%s: %s\n", state, targetURL)
		} else {
			ui.Println(state)
		}

		os.Exit(exitCode)
	}
}

func fetchCiStatus(p *github.Project, sha string) (state, targetURL string, exitCode int, err error) {
	gh := github.NewClient(p.Host)
	status, err := gh.CIStatus(p, sha)
	if err != nil {
		return
	}

	if status == nil {
		state = "no status"
	} else {
		state = status.State
		targetURL = status.TargetURL
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
