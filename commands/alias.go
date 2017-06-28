package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var cmdAlias = &Command{
	Run:   alias,
	Usage: "alias [-s] [<SHELL>]",
	Long: `Show shell instructions for wrapping git.

## Options
	-s
		Output shell script suitable for 'eval'.

	<SHELL>
		Specify the type of shell (default: "$SHELL" environment variable).

## See also:

hub(1)
`,
}

var flagAliasScript bool

func init() {
	cmdAlias.Flag.BoolVarP(&flagAliasScript, "script", "s", false, "SCRIPT")
	CmdRunner.Use(cmdAlias)
}

func alias(command *Command, args *Args) {
	var shell string
	if args.ParamsSize() > 0 {
		shell = args.FirstParam()
	} else {
		shell = os.Getenv("SHELL")
	}

	if shell == "" {
		cmd := "hub alias <shell>"
		if flagAliasScript {
			cmd = "hub alias -s <shell>"
		}
		utils.Check(fmt.Errorf("Error: couldn't detect shell type. Please specify your shell with `%s`", cmd))
	}

	shells := []string{"bash", "zsh", "sh", "ksh", "csh", "tcsh", "fish", "rc"}
	shell = filepath.Base(shell)
	var validShell bool
	for _, s := range shells {
		if s == shell {
			validShell = true
			break
		}
	}

	if !validShell {
		err := fmt.Errorf("hub alias: unsupported shell\nsupported shells: %s", strings.Join(shells, " "))
		utils.Check(err)
	}

	if flagAliasScript {
		var alias string
		switch shell {
		case "csh", "tcsh":
			alias = "alias git hub"
		case "rc":
			alias = "fn git { builtin hub $* }"
		default:
			alias = "alias git=hub"
		}

		ui.Println(alias)
	} else {
		var profile string
		switch shell {
		case "bash":
			profile = "~/.bash_profile"
		case "zsh":
			profile = "~/.zshrc"
		case "ksh":
			profile = "~/.profile"
		case "fish":
			profile = "~/.config/fish/functions/git.fish"
		case "csh":
			profile = "~/.cshrc"
		case "tcsh":
			profile = "~/.tcshrc"
		case "rc":
			profile = "$home/lib/profile"
		default:
			profile = "your profile"
		}

		msg := fmt.Sprintf("# Wrap git automatically by adding the following to %s:\n", profile)
		ui.Println(msg)

		var eval string
		switch shell {
		case "fish":
			eval = `function git --wraps hub --description 'Alias for hub, which wraps git to provide extra functionality with GitHub.'
	hub $argv
end`
		case "rc":
			eval = "eval `{hub alias -s}"
		case "csh", "tcsh":
			eval = "eval \"`hub alias -s`\""
		default:
			eval = `eval "$(hub alias -s)"`
		}

		indent := regexp.MustCompile(`(?m)^\t+`)
		eval = indent.ReplaceAllStringFunc(eval, func(match string) string {
			return strings.Repeat(" ", len(match)*4)
		})

		ui.Println(eval)
	}

	args.NoForward()
}
