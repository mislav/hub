package commands

import (
	"encoding/json"
	"fmt"

	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var (
	cmdGist = &Command{
		Run: printGistHelp,
		Usage: `
gist create [--public] [<FILES>...]
gist show [--json] <ID> [<FILENAME>]
`,
		Long: `Create and print GitHub Gists

## Commands:

	* _create_:
		Create a new gist. If no <FILES> are specified, the content is read from
		standard input.

    * _show_:
		Print the contents of a gist. If the gist contains multiple files, the
		operation will error out unless either <FILENAME> is specified or the
		'--json' flag is used.

## Options:

    --public
        Make the new gist public (default: false).

    --json
        Print all files in the gist and emit them in JSON.

## Examples:

	$ echo hello | hub gist create --public

	$ hub gist create <file1> <file2>

    # print a specific file within a gist:
    $ hub gist show <ID> testfile1.txt

## See also:

hub(1), hub-api(1)
`,
		KnownFlags: "\n",
	}

	cmdShowGist = &Command{
		Key:        "show",
		Run:        showGist,
		KnownFlags: "--json",
	}

	cmdCreateGist = &Command{
		Key:        "create",
		Run:        createGist,
		KnownFlags: "--public",
	}
)

func init() {
	cmdGist.Use(cmdShowGist)
	cmdGist.Use(cmdCreateGist)
	CmdRunner.Use(cmdGist)
}

func getGist(gh *github.Client, id string, filename string, emitJson bool) error {
	gist, err := gh.FetchGist(id)
	if err != nil {
		return err
	}

	if len(gist.Files) > 1 && !emitJson && filename == "" {
		return fmt.Errorf("There are multiple files, you must specify one, or use --json")
	}

	if emitJson {
		data, err := json.Marshal(gist.Files)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", data)
	} else if filename != "" {
		if val, ok := gist.Files[filename]; ok {
			ui.Println(val.Content)
		} else {
			return fmt.Errorf("no such file in gist")
		}
	} else {
		/*
		 * There's only one, but we don't know the name so a
		 * loop us fine.
		 */
		for name := range gist.Files {
			file := gist.Files[name]
			ui.Println(file.Content)
		}
	}
	return nil
}

func printGistHelp(command *Command, args *Args) {
	utils.Check(command.UsageError(""))
}

func createGist(cmd *Command, args *Args) {
	args.NoForward()

	host, err := github.CurrentConfig().DefaultHostNoPrompt()
	utils.Check(err)
	gh := github.NewClient(host.Host)

	filenames := []string{}
	if args.IsParamsEmpty() {
		filenames = append(filenames, "-")
	} else {
		filenames = args.Params
	}
	g, err := gh.CreateGist(filenames, args.Flag.Bool("--public"))
	utils.Check(err)
	ui.Println(g.HtmlUrl)
}

func showGist(cmd *Command, args *Args) {
	args.NoForward()

	host, err := github.CurrentConfig().DefaultHostNoPrompt()
	utils.Check(err)
	gh := github.NewClient(host.Host)

	id := args.GetParam(0)
	filename := ""
	if args.ParamsSize() > 1 {
		filename = args.GetParam(1)
	}
	err = getGist(gh, id, filename, args.Flag.Bool("--json"))
	utils.Check(err)
}
