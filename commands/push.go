package commands

import (
	"fmt"
	"os"
	"strings"
	"github.com/jingweno/gh/utils"
)

var cmdPush = &Command{
	Run: push,
	GitExtension: true,
	Usage: "push REMOTE-1,REMOTE-2,...,REMOTE-N [REF]",
	Short: "Update remote refs along with associated objects",
	Long: `Push REF to each of REMOTE-1 through REMOTE-N by executing  mul-
tiple git push commands.`,
}

/**
 $ git push origin,staging,qa bert_timeout
 > git push origin bert_timeout
 > git push staging bert_timeout
 > git push qa bert_timeout
**/

func push (command *Command, args *Args) {
	if !args.IsParamsEmpty() {
		pushToEveryRemote(args)
	}
}

func pushToEveryRemote (args *Args) {
	remotes, idx := getRemotes(args)
	for _, i := range remotes {
		copyArgs := args
		copyArgs.ReplaceParam(idx, i)
		if !args.Noop {
			err := copyArgs.ToCmd().Exec()
			utils.Check(err)
		} else {
			fmt.Printf("it would run `git push %s`\n", strings.Join(copyArgs.Params, " "))
		}
	}

	fixHelp(args)
}

func getRemotes(args *Args) (remotes []string, idx int) {
	for a, i := range args.Params {
		if !strings.HasPrefix(i, "-") {
			remotes = strings.Split(i, ",")
			idx = a
			return
		}
	}

	return
}

func fixHelp(args *Args) {
	if i := args.IndexOfParam("--help"); i != -1 && args.ParamsSize() == 1 {
		args.Params = []string{"--help"}
		return
	}

	os.Exit(0)
}
