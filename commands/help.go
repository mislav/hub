package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/github/hub/cmd"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var cmdHelp = &Command{
	Run:          runHelp,
	GitExtension: true,
	Usage: `
help hub
help <COMMAND>
help hub-<COMMAND> [--plain-text]
`,
	Long: `Show the help page for a command.

## Options:
	hub-<COMMAND>
		Use this format to view help for hub extensions to an existing git command.

	--plain-text
		Skip man page lookup mechanism and display plain help text.

## Man lookup mechanism:

On systems that have 'man', help pages are looked up in these directories
relative to 'hub' install prefix:

* 'man/<command>.1'
* 'share/man/man1/<command>.1'

On systems without 'man', same help pages are looked up with a '.txt' suffix.

## See also:

hub(1), git-help(1)
`,
}

func init() {
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

	if args.HasFlags("-a", "--all") {
		args.AfterFn(func() error {
			ui.Printf("\nhub custom commands\n\n  %s\n", strings.Join(customCommands(), "  "))
			return nil
		})
		return
	}

	command := args.FirstParam()

	if command == "hub" {
		err := displayManPage("hub.1", args)
		if err != nil {
			utils.Check(err)
		}
	}

	if c := lookupCmd(command); c != nil {
		if !args.HasFlags("--plain-text") {
			manPage := fmt.Sprintf("hub-%s.1", c.Name())
			err := displayManPage(manPage, args)
			if err == nil {
				return
			}
		}

		ui.Println(c.HelpText())
		args.NoForward()
	}
}

func displayManPage(manPage string, args *Args) error {
	manProgram, _ := utils.CommandPath("man")
	if manProgram == "" {
		manPage += ".txt"
		manProgram = os.Getenv("PAGER")
		if manProgram == "" {
			manProgram = "less -R"
		}
	}

	programPath, err := utils.CommandPath(args.ProgramPath)
	if err != nil {
		return err
	}

	installPrefix := filepath.Join(filepath.Dir(programPath), "..")
	manFile, err := localManPage(manPage, installPrefix)
	if err != nil {
		return err
	}

	man := cmd.New(manProgram)
	man.WithArg(manFile)
	if err = man.Run(); err == nil {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
	return nil
}

func localManPage(name, installPrefix string) (string, error) {
	manPath := filepath.Join(installPrefix, "man", name)
	_, err := os.Stat(manPath)
	if err == nil {
		return manPath, nil
	}

	manPath = filepath.Join(installPrefix, "share", "man", "man1", name)
	_, err = os.Stat(manPath)
	if err == nil {
		return manPath, nil
	} else {
		return "", err
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

   browse         Open a GitHub page in the default browser
   ci-status      Show the CI status of a commit
   compare        Open a compare page on GitHub
   create         Create this repository on GitHub and add GitHub as origin
   fork           Make a fork of a remote repository on GitHub and add as remote
   issue          List or create issues
   pr             Work with pull requests
   pull-request   Open a pull request on GitHub
   release        List or create releases
`
