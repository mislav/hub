package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/github/hub/git"
	"github.com/github/hub/utils"
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

func runHelp(cmd *Command, args *Args) {
	if args.IsParamsEmpty() {
		printUsage()
		os.Exit(0)
	}

	command := args.FirstParam()
	c := CmdRunner.Lookup(command)
	if c != nil && !c.GitExtension {
		c.PrintUsage()
		os.Exit(0)
	} else if c == nil {
		if args.HasFlags("-a", "--all") {
			args.After("echo", "\nhub custom commands\n")
			args.After("echo", " ", strings.Join(customCommands(), "  "))
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

func printUsage() {
	err := git.ForwardGitHelp()
	utils.Check(err)
	fmt.Print(helpText)
}
