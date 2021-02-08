package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/github/hub/v2/git"
	"github.com/github/hub/v2/ui"
	"github.com/github/hub/v2/utils"
	"github.com/kballard/go-shellquote"
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
		Skip man page lookup mechanism and display raw help text.

## See also:

hub(1), git-help(1)
`,
}

var cmdListCmds = &Command{
	Key:          "--list-cmds",
	Run:          runListCmds,
	GitExtension: true,
}

func init() {
	CmdRunner.Use(cmdHelp, "--help")
	CmdRunner.Use(cmdListCmds)
}

func runHelp(helpCmd *Command, args *Args) {
	if args.IsParamsEmpty() {
		args.AfterFn(func() error {
			ui.Println(helpText)
			return nil
		})
		return
	}

	p := utils.NewArgsParser()
	p.RegisterBool("--all", "-a")
	p.RegisterBool("--plain-text")
	p.RegisterBool("--man", "-m")
	p.RegisterBool("--web", "-w")
	p.Parse(args.Params)

	if p.Bool("--all") {
		args.AfterFn(func() error {
			ui.Printf("\nhub custom commands\n\n  %s\n", strings.Join(customCommands(), "  "))
			return nil
		})
		return
	}

	isWeb := func() bool {
		if p.Bool("--web") {
			return true
		}
		if p.Bool("--man") {
			return false
		}
		if f, err := git.Config("help.format"); err == nil {
			return f == "web" || f == "html"
		}
		return false
	}

	cmdName := ""
	if words := args.Words(); len(words) > 0 {
		cmdName = words[0]
	}

	if cmdName == "hub" {
		err := displayManPage("hub", args, isWeb())
		utils.Check(err)
		return
	}

	foundCmd := lookupCmd(cmdName)
	if foundCmd == nil {
		return
	}

	if p.Bool("--plain-text") {
		ui.Println(foundCmd.HelpText())
		os.Exit(0)
	}

	manPage := fmt.Sprintf("hub-%s", foundCmd.Name())
	err := displayManPage(manPage, args, isWeb())
	utils.Check(err)
}

func runListCmds(cmd *Command, args *Args) {
	listOthers := false
	parts := strings.SplitN(args.Command, "=", 2)
	for _, kind := range strings.Split(parts[1], ",") {
		if kind == "others" {
			listOthers = true
			break
		}
	}

	if listOthers {
		args.AfterFn(func() error {
			ui.Println(strings.Join(customCommands(), "\n"))
			return nil
		})
	}
}

// On systems where `man` was found, invoke:
//   MANPATH={PREFIX}/share/man:$MANPATH man <page>
//
// otherwise:
//   less -R {PREFIX}/share/man/man1/<page>.1.txt
func displayManPage(manPage string, args *Args, isWeb bool) error {
	programPath, err := utils.CommandPath(args.ProgramPath)
	if err != nil {
		return err
	}

	if isWeb {
		manPage += ".1.html"
		manFile := filepath.Join(programPath, "..", "..", "share", "doc", "hub-doc", manPage)
		args.Replace(args.Executable, "web--browse", manFile)
		return nil
	}

	var manArgs []string
	manProgram, _ := utils.CommandPath("man")
	if manProgram != "" {
		manArgs = []string{manProgram}
	} else {
		manPage += ".1.txt"
		if manProgram = os.Getenv("PAGER"); manProgram != "" {
			var err error
			manArgs, err = shellquote.Split(manProgram)
			if err != nil {
				return err
			}
		} else {
			manArgs = []string{"less", "-R"}
		}
	}

	env := os.Environ()
	if strings.HasSuffix(manPage, ".txt") {
		manFile := filepath.Join(programPath, "..", "..", "share", "man", "man1", manPage)
		manArgs = append(manArgs, manFile)
	} else {
		manArgs = append(manArgs, manPage)
		manPath := filepath.Join(programPath, "..", "..", "share", "man")
		env = append(env, fmt.Sprintf("MANPATH=%s:%s", manPath, os.Getenv("MANPATH")))
	}

	c := exec.Command(manArgs[0], manArgs[1:]...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = env
	if err := c.Run(); err != nil {
		return err
	}
	os.Exit(0)
	return nil
}

func lookupCmd(name string) *Command {
	if strings.HasPrefix(name, "hub-") {
		return CmdRunner.Lookup(strings.TrimPrefix(name, "hub-"))
	}
	cmd := CmdRunner.Lookup(name)
	if cmd != nil && !cmd.GitExtension {
		return cmd
	}
	return nil
}

func customCommands() []string {
	cmds := []string{}
	for n, c := range CmdRunner.All() {
		if !c.GitExtension && !strings.HasPrefix(n, "--") {
			cmds = append(cmds, n)
		}
	}

	sort.Strings(cmds)

	return cmds
}

var helpText = `
These GitHub commands are provided by hub:

   api            Low-level GitHub API request interface
   browse         Open a GitHub page in the default browser
   ci-status      Show the status of GitHub checks for a commit
   compare        Open a compare page on GitHub
   create         Create this repository on GitHub and add GitHub as origin
   delete         Delete a repository on GitHub
   fork           Make a fork of a remote repository on GitHub and add as remote
   gist           Make a gist
   issue          List or create GitHub issues
   pr             Manage GitHub pull requests
   pull-request   Open a pull request on GitHub
   release        List or create GitHub releases
   sync           Fetch git objects from upstream and update branches
`
