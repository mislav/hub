package commands

import (
	"fmt"
	"github.com/github/hub/utils"
	"os"
	"path/filepath"
	"strings"
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

	shells := []string{"bash", "zsh", "sh", "ksh", "csh", "fish"}
	shell = filepath.Base(shell)
	var validShell bool
	for _, s := range shells {
		if s == shell {
			validShell = true
			break
		}
	}

	if !validShell {
		err := fmt.Errorf("gh alias: unsupported shell\nsupported shells: %s", strings.Join(shells, " "))
		utils.Check(err)
	}

	if flagAliasScript {
		fmt.Println("alias git=gh")
		if "zsh" == shell {
			fmt.Println("if type compdef > /dev/null; then")
			fmt.Println("  compdef gh=git")
			fmt.Println("fi")
		}
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
		default:
			profile = "your profile"
		}

		msg := fmt.Sprintf("# Wrap git automatically by adding the following to %s:\n", profile)
		fmt.Println(msg)

		var eval string
		if shell == "fish" {
			eval = `eval (gh alias -s)`
		} else {
			eval = `eval "$(gh alias -s)"`
		}
		fmt.Println(eval)
	}

	os.Exit(0)
}
