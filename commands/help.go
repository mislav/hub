package commands

import (
	"fmt"
	"github.com/jingweno/gh/utils"
	"os"
	"text/template"
)

var cmdHelp = &Command{
	Usage: "help [command]",
	Short: "Show help",
	Long:  `Shows usage for a command.`,
}

func init() {
	cmdHelp.Run = runHelp // break init loop
}

func runHelp(cmd *Command, args *Args) {
	if args.IsParamsEmpty() {
		printUsage()
		os.Exit(0)
	}

	if args.ParamsSize() > 1 {
		utils.Check(fmt.Errorf("too many arguments"))
	}

	for _, cmd := range All() {
		if cmd.Name() == args.FirstParam() {
			cmd.PrintUsage()
			os.Exit(0)
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic: %q. Run 'gh help'.\n", args.FirstParam())
	os.Exit(2)
}

var usageTemplate = template.Must(template.New("usage").Parse(`Usage: gh [command] [options] [arguments]

Branching Commands:{{range .BasicCommands}}{{if .Runnable}}{{if .List}}
    {{.Name | printf "%-16s"}}  {{.Short}}{{end}}{{end}}{{end}}

Branching Commands:{{range .BranchingCommands}}{{if .Runnable}}{{if .List}}
    {{.Name | printf "%-16s"}}  {{.Short}}{{end}}{{end}}{{end}}

Remote Commands:{{range .RemoteCommands}}{{if .Runnable}}{{if .List}}
    {{.Name | printf "%-16s"}}  {{.Short}}{{end}}{{end}}{{end}}

GitHub Commands:{{range .GitHubCommands}}{{if .Runnable}}{{if .List}}
    {{.Name | printf "%-16s"}}  {{.Short}}{{end}}{{end}}{{end}}

See 'gh help [command]' for more information about a command.
`))

func printUsage() {
	usageTemplate.Execute(os.Stdout, struct {
		BasicCommands     []*Command
		BranchingCommands []*Command
		RemoteCommands    []*Command
		GitHubCommands    []*Command
	}{
		Basic,
		Branching,
		Remote,
		GitHub,
	})
}

func usage() {
	printUsage()
	os.Exit(2)
}
