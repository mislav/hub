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
gist show <GistID> [--json] [<filename>]
gist create [--public] <file1> [<file2> ..]
cat <file> | gist create [--public]
`,
		Long: `Create a GitHub Gist

Can both create and retrieve gists. With no arguements, it takes a file on
stdin and creates a gist. With multiple '--file' arguments, will create a
mult-file gist.

If gistid is passed in, if there is only one file in the gist, it will be
printed, otherwise will error if a specific file is not requested. However, if
'--json' is used, then multiple files can be retreived.

## Commands:

    * _show_:
        Show the gist <GistID>. If the gist has more than one file in it,
        either a file must be specified, or --json must be used. When --json
        is specified, filenames are ignored.

    * _create_:
        Create a gist. If no files are specified, the file is content
        is read from stdin. You may specify as many files as you want.

## Options:

    --public
        The gist should be marked as public.

    --json
        Print all files in the gist and emit them in JSON.

## Examples:

    # Retrieve the contents of a gist with a single file
    $ hub gist show 87560fa4ebcc8683f68ec04d9100ab1c
    this is test content in testfile.txt in test gist

    # Retrieve same gist, but specify a single file
    $ hub gist show 6188fb16b1a7df0f51a51e03b3a2b4e8 testfile1.txt
    test content in testfile1.txt inside of test gist 2

    # Retrieve same gist, with all files, using JSON
    $ hub gist show --json 6188fb16b1a7df0f51a51e03b3a2b4e8
    test content in testfile1.txt inside of test gist 2
    more test content in testfile2.txt inside of test gist 2

    # Create a gist:
    $ cat /tmp/testfile | hub gist create
    https://gist.github.com/bdf551042f77bb8431b99f13c1105168

    # Or a public one:
    $ cat /tmp/testfile | hub gist create --public
    https://gist.github.com/6c925133a295f0c5ad61eafcf05fee30

    # You can also specify a file directly
    $ hub gist create /tmp/testfile
    https://gist.github.com/bdf551042f77bb8431b99f13c1105168

    # Or with multiple files
    $ hub gist create /tmp/testfile /tmp/testfile2
    https://gist.github.com/bdf551042f77bb8431b99f13c1105168


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

	host, err := github.CurrentConfig().DefaultHost()
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

	host, err := github.CurrentConfig().DefaultHost()
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
