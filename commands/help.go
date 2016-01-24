package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/github/hub/cmd"
	"github.com/github/hub/git"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var cmdHelp = &Command{
	Run:          runHelp,
	GitExtension: true,
	Usage:        "help [<COMMAND>]",
	Long:         `Show the help page for a command.`,
}

func init() {
	CmdRunner.Use(cmdHelp, "--help")
}

func runHelp(helpCmd *Command, args *Args) {
	if args.IsParamsEmpty() {
		printUsage()
		os.Exit(0)
	}

	if args.HasFlags("-a", "--all") {
		args.After("echo", "\nhub custom commands\n")
		args.After("echo", " ", strings.Join(customCommands(), "  "))
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

	if c := lookupCmd(command); c != nil {
		ui.Println(c.HelpText())
		os.Exit(0)
	}
}

func lookupCmd(name string) *Command {
	if strings.HasPrefix(name, "hub-") {
		return CmdRunner.Lookup(strings.TrimPrefix(name, "hub-"))
	} else {
		cmd := CmdRunner.Lookup(name)
		if cmd != nil && !cmd.GitExtension {
			return cmd
		} else {
			return nil
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
