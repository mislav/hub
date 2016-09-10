package commands

import (
	"os"
	"sort"
	"strings"

	"github.com/github/hub/cmd"
	"github.com/github/hub/ui"
)

var cmdHelp = &Command{
	Usage:        "help [command]",
	Short:        "Show help",
	Long:         `Shows usage for a command.`,
	GitExtension: true,
}

func init() {
	cmdHelp.Run = runHelp

	CmdRunner.Use(cmdHelp, "--help")
}

func runHelp(helpCmd *Command, args *Args) {
	if args.IsParamsEmpty() {
		args.AfterFn(func() error {
			ui.Println(helpText)
			return nil
		})
		return
	}

	command := args.FirstParam()

	if command == "hub" {
		man := cmd.New("man")
		man.WithArg("hub")
		err := man.Run()
		if err == nil {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	c := CmdRunner.Lookup(command)
	if c != nil && !c.GitExtension {
		c.PrintUsage()
		args.NoForward()
	} else if c == nil {
		if args.HasFlags("-a", "--all") {
			args.AfterFn(func() error {
				ui.Printf("\nhub custom commands\n\n  %s\n", strings.Join(customCommands(), "  "))
				return nil
			})
		}
	}
}

func customCommands() []string {
	cmds := []string{}
	for n, c := range CmdRunner.All() {
		if !c.GitExtension && !strings.HasPrefix(n, "--") {
			cmds = append(cmds, n)
		}
	}

	sort.Sort(sort.StringSlice(cmds))

	return cmds
}

var helpText = `
These GitHub commands are provided by hub:

   pull-request   Open a pull request on GitHub
   fork           Make a fork of a remote repository on GitHub and add as remote
   create         Create this repository on GitHub and add GitHub as origin
   browse         Open a GitHub page in the default browser
   compare        Open a compare page on GitHub
   release        List or create releases (beta)
   issue          List or create issues (beta)
   ci-status      Show the CI status of a commit
`
