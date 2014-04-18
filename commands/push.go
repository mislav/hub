package commands

import (
	"strings"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var cmdPush = &Command{
	Run:          push,
	GitExtension: true,
	Usage:        "push REMOTE-1,REMOTE-2,...,REMOTE-N [REF]",
	Short:        "Upload data, tags and branches to a remote repository",
	Long: `Push REF to each of REMOTE-1 through REMOTE-N by executing
multiple git-push(1) commands.`,
}

func init() {
	CmdRunner.Use(cmdPush)
}

/*
  $ gh push origin,staging,qa bert_timeout
  > git push origin bert_timeout
  > git push staging bert_timeout
  > git push qa bert_timeout

  $ gh push origin
  > git push origin HEAD
*/
func push(command *Command, args *Args) {
	if !args.IsParamsEmpty() && strings.Contains(args.FirstParam(), ",") {
		transformPushArgs(args)
	}
}

func transformPushArgs(args *Args) {
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
