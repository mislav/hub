package commands

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/github/hub/git"
	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var cmdSync = &Command{
	Run:   sync,
	Usage: "sync",
	Long: `Fetch git objects from upstream and update branches.

- If the local branch is outdated, fast-forward it;
- If the local branch contains unpushed work, warn about it;
- If the branch seems merged and its upstream branch was deleted, delete it.

If a local branch doesn't have any upstream configuration, but has a
same-named branch on the remote, treat that as its upstream branch.

## See also:

hub(1), git-fetch(1)
`,
}

func init() {
	CmdRunner.Use(cmdSync)
}

func sync(cmd *Command, args *Args) {
	localRepo, err := github.LocalRepo()
	utils.Check(err)

	remote, err := localRepo.MainRemote()
	utils.Check(err)

	defaultBranch := localRepo.MasterBranch().ShortName()
	fullDefaultBranch := fmt.Sprintf("refs/remotes/%s/%s", remote.Name, defaultBranch)
	currentBranch := ""
	if curBranch, err := localRepo.CurrentBranch(); err == nil {
		currentBranch = curBranch.ShortName()
	}

	err = git.Spawn("fetch", "--prune", "--quiet", "--progress", remote.Name)
	utils.Check(err)

	branchToRemote := map[string]string{}
	if lines, err := git.ConfigAll("branch.*.remote"); err == nil {
		configRe := regexp.MustCompile(`^branch\.(.+?)\.remote (.+)`)

		for _, line := range lines {
			if matches := configRe.FindStringSubmatch(line); len(matches) > 0 {
				branchToRemote[matches[1]] = matches[2]
			}
		}
	}

	branches, err := git.LocalBranches()
	utils.Check(err)

	var green,
		lightGreen,
		red,
		lightRed,
		resetColor string

	if ui.IsTerminal(os.Stdout) {
		green = "\033[32m"
		lightGreen = "\033[32;1m"
		red = "\033[31m"
		lightRed = "\033[31;1m"
		resetColor = "\033[0m"
	}

	for _, branch := range branches {
		fullBranch := fmt.Sprintf("refs/heads/%s", branch)
		remoteBranch := fmt.Sprintf("refs/remotes/%s/%s", remote.Name, branch)
		gone := false

		if branchToRemote[branch] == remote.Name {
			if upstream, err := git.SymbolicFullName(fmt.Sprintf("%s@{upstream}", branch)); err == nil {
				remoteBranch = upstream
			} else {
				remoteBranch = ""
				gone = true
			}
		} else if !git.HasFile(strings.Split(remoteBranch, "/")...) {
			remoteBranch = ""
		}

		if remoteBranch != "" {
			diff, err := git.NewRange(fullBranch, remoteBranch)
			utils.Check(err)

			if diff.IsIdentical() {
				continue
			} else if diff.IsAncestor() {
				if branch == currentBranch {
					git.Quiet("merge", "--ff-only", "--quiet", remoteBranch)
				} else {
					git.Quiet("update-ref", fullBranch, remoteBranch)
				}
				ui.Printf("%sUpdated branch %s%s%s (was %s).\n", green, lightGreen, branch, resetColor, diff.A[0:7])
			} else {
				ui.Errorf("warning: `%s' seems to contain unpushed commits\n", branch)
			}
		} else if gone {
			diff, err := git.NewRange(fullBranch, fullDefaultBranch)
			utils.Check(err)

			if diff.IsAncestor() {
				if branch == currentBranch {
					git.Quiet("checkout", "--quiet", defaultBranch)
					currentBranch = defaultBranch
				}
				git.Quiet("branch", "-D", branch)
				ui.Printf("%sDeleted branch %s%s%s (was %s).\n", red, lightRed, branch, resetColor, diff.A[0:7])
			} else {
				ui.Errorf("warning: `%s' was deleted on %s, but appears not merged into %s\n", branch, remote.Name, defaultBranch)
			}
		}
	}

	args.NoForward()
}
