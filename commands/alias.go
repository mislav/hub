package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var cmdAlias = &Command{
	Run:   alias,
	Usage: "alias [-s] [SHELL]",
	Short: "Show shell instructions for wrapping git",
	Long: `Shows shell instructions for wrapping git. If given, SHELL specifies the
type of shell; otherwise defaults to the value of SHELL environment
variable. With -s, outputs shell script suitable for eval.
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
		utils.Check(fmt.Errorf("Unknown shell"))
	}

	shells := []string{"bash", "zsh", "sh", "ksh", "csh", "tcsh", "fish"}
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
			profile = "~/.config/fish/config.fish"
		case "csh":
			profile = "~/.cshrc"
		case "tcsh":
			profile = "~/.tcshrc"
		default:
			profile = "your profile"
		}

		msg := fmt.Sprintf("# Wrap git automatically by adding the following to %s:\n", profile)
		ui.Println(msg)

		var eval string
		switch shell {
		case "fish":
			eval = `eval (hub alias -s)`
		case "csh", "tcsh":
			eval = "eval \"`hub alias -s`\""
		default:
			eval = `eval "$(hub alias -s)"`
		}
		ui.Println(eval)
	}

	os.Exit(0)
}
