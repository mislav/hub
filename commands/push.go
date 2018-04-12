package commands

import (
  "fmt"
	"strings"

	"github.com/github/hub/git"
	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var cmdPush = &Command{
	Run:          push,
	GitExtension: true,
	Usage:        "push <REMOTE>[,<REMOTE2>...] [<REF>]",
	Long: `Push a git branch to each of the listed remotes.

## Examples:
		$ hub push origin,staging,qa bert_timeout
		> git push origin bert_timeout
		> git push staging bert_timeout
		> git push qa bert_timeout

		$ hub push origin
		> git push origin HEAD

    TODO: something about pushBranch

## See also:

hub(1), git-push(1)
`,
}

func init() {
	CmdRunner.Use(cmdPush)
}

func push(command *Command, args *Args) {
	if args.IsParamsEmpty() {
    pushRemoteBranch(args)
  } else if strings.Contains(args.FirstParam(), ",") {
		pushMultipleRemotes(args)
	}
}

func pushRemoteBranch(args *Args) {
  localRepo, err := github.LocalRepo()
  utils.Check(err)

  currentBranch, err := localRepo.CurrentBranch()
  utils.Check(err)

  shortName := currentBranch.ShortName()
  pushRemote, _ := git.Config(fmt.Sprintf("branch.%s.pushRemote", shortName))
  if pushRemote != "" {
    pushBranch, _ := git.Config(fmt.Sprintf("branch.%s.pushBranch", shortName))
    if pushBranch != "" {
      args.AppendParams(pushRemote, fmt.Sprintf("%s:%s", shortName, pushBranch))
    }
  }
}

func pushMultipleRemotes(args *Args) {
	refs := []string{}
	if args.ParamsSize() > 1 {
		refs = args.Params[1:]
	}

	remotes := strings.Split(args.FirstParam(), ",")
	args.ReplaceParam(0, remotes[0])

	if len(refs) == 0 {
		localRepo, err := github.LocalRepo()
		utils.Check(err)

		head, err := localRepo.CurrentBranch()
		utils.Check(err)

		refs = []string{head.ShortName()}
		args.AppendParams(refs...)
	}

	for _, remote := range remotes[1:] {
		afterCmd := []string{"git", "push", remote}
		afterCmd = append(afterCmd, refs...)
		args.After(afterCmd...)
	}
}
