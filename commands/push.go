package commands

import (
	"github.com/jingweno/gh/github"
	"github.com/jingweno/gh/utils"
	"strings"
)

var cmdPush = &Command{
	Run:          push,
	GitExtension: true,
	Usage:        "push REMOTE-1,REMOTE-2,...,REMOTE-N [REF]",
	Short:        "Upload data, tags and branches to a remote repository",
	Long: `Push REF to each of REMOTE-1 through REMOTE-N by executing
multiple git-push(1) commands.`,
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
		localRepo := github.LocalRepo()
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
