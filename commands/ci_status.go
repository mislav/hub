package commands

import (
	//"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/utils"
)

var cmdCiStatus = &Command{
	Run:   ciStatus,
	Usage: "ci-status [COMMIT]",
	Short: "Show CI status of a commit",
	Long: `Looks up the SHA for COMMIT in GitHub Status API and displays the latest
status. Exits with one of:

success (0), error (1), failure (1), pending (2), no status (3)
`,
}

func ciStatus(cmd *Command, args []string) {
	if len(args) == 0 {
		cmd.PrintUsage()
		return
	}

	ref := args[0]
	if ref == "" {
		ref = "HEAD"
	}

	ref, err := git.Ref(ref)
	utils.Check(err)

	//github := github.NewGitHub()
	//github.ListStatuses(github.CurrentProject(), ref)
}
