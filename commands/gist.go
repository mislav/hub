package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/github/hub/v2/github"
	"github.com/github/hub/v2/ui"
	"github.com/github/hub/v2/utils"
)

var (
	cmdGist = &Command{
		Run: printGistHelp,
		Usage: `
gist create [-oc] [--public] [<FILES>...]
gist show <ID> [<FILENAME>]
`,
		Long: `Create and print GitHub Gists

## Commands:

	* _create_:
		Create a new gist. If no <FILES> are specified, the content is read from
		standard input.

	* _show_:
		Print the contents of a gist. If the gist contains multiple files, the
		operation will error out unless <FILENAME> is specified.

## Options:

	--public
		Make the new gist public (default: false).

	-o, --browse
		Open the new gist in a web browser.

	-c, --copy
		Put the URL of the new gist to clipboard instead of printing it.

## Examples:

    $ echo hello | hub gist create --public

    $ hub gist create file1.txt file2.txt

    # print a specific file within a gist:
    $ hub gist show ID testfile1.txt

## See also:

hub(1), hub-api(1)
`,
	}

	cmdShowGist = &Command{
		Key: "show",
		Run: showGist,
	}

	cmdCreateGist = &Command{
		Key: "create",
		Run: createGist,
		KnownFlags: `
		--public
		-o, --browse
		-c, --copy
`,
	}
)

func init() {
	cmdGist.Use(cmdShowGist)
	cmdGist.Use(cmdCreateGist)
	CmdRunner.Use(cmdGist)
}

func getGist(gh *github.Client, id string, filename string) error {
	gist, err := gh.FetchGist(id)
	if err != nil {
		return err
	}

	if len(gist.Files) > 1 && filename == "" {
		filenames := []string{}
		for name := range gist.Files {
			filenames = append(filenames, name)
		}
		sort.Strings(filenames)
		return fmt.Errorf("This gist contains multiple files, you must specify one:\n  %s", strings.Join(filenames, "\n  "))
	}

	if filename != "" {
		if val, ok := gist.Files[filename]; ok {
			ui.Println(val.Content)
		} else {
			return fmt.Errorf("no such file in gist")
		}
	} else {
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

	var gist *github.Gist
	if args.Noop {
		ui.Println("Would create gist")
		gist = &github.Gist{
			HTMLURL: fmt.Sprintf("https://gist.%s/%s", gh.Host.Host, "ID"),
		}
	} else {
		gist, err = gh.CreateGist(filenames, args.Flag.Bool("--public"))
		utils.Check(err)
	}

	flagIssueBrowse := args.Flag.Bool("--browse")
	flagIssueCopy := args.Flag.Bool("--copy")
	printBrowseOrCopy(args, gist.HTMLURL, flagIssueBrowse, flagIssueCopy)
}

func showGist(cmd *Command, args *Args) {
	args.NoForward()

	if args.ParamsSize() < 1 {
		utils.Check(cmd.UsageError("you must specify a gist ID"))
	}

	host, err := github.CurrentConfig().DefaultHostNoPrompt()
	utils.Check(err)
	gh := github.NewClient(host.Host)

	id := args.GetParam(0)
	filename := ""
	if args.ParamsSize() > 1 {
		filename = args.GetParam(1)
	}

	err = getGist(gh, id, filename)
	utils.Check(err)
}
