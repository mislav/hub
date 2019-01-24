package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var cmdApi = &Command{
	Run:   apiCommand,
	Usage: "api <RESOURCE>",
	Long: `Interact with the GitHub API.

## Options:
	-X, --method <METHOD>
	-F, --field <KEY-VALUE>
	-f, --raw-field <KEY-VALUE>
	-t, --flat
	--cache <TTL>
`,
}

func init() {
	CmdRunner.Use(cmdApi)
}

func apiCommand(cmd *Command, args *Args) {
	path := ""
	if !args.IsParamsEmpty() {
		path = args.GetParam(0)
	}

	method := "GET"
	if args.Flag.HasReceived("--method") {
		method = args.Flag.Value("--method")
	} else if args.Flag.HasReceived("--field") {
		method = "POST"
	}
	cacheTTL := args.Flag.Int("--cache")

	params := make(map[string]interface{})
	for _, val := range args.Flag.AllValues("--field") {
		parts := strings.SplitN(val, "=", 2)
		if len(parts) >= 2 {
			params[parts[0]] = valueOrFileContents(parts[1])
		}
	}
	for _, val := range args.Flag.AllValues("--raw-field") {
		parts := strings.SplitN(val, "=", 2)
		if len(parts) >= 2 {
			params[parts[0]] = parts[1]
		}
	}

	host := ""
	owner := ""
	repo := ""
	localRepo, localRepoErr := github.LocalRepo()
	if localRepoErr == nil {
		var project *github.Project
		if project, localRepoErr = localRepo.MainProject(); localRepoErr == nil {
			host = project.Host
			owner = project.Owner
			repo = project.Name
		}
	}
	if host == "" {
		defHost, err := github.CurrentConfig().DefaultHost()
		utils.Check(err)
		host = defHost.Host
	}

	if path == "graphql" && params["query"] != nil {
		query := params["query"].(string)
		query = strings.Replace(query, quote("{owner}"), quote(owner), 1)
		query = strings.Replace(query, quote("{repo}"), quote(repo), 1)
		params["query"] = query
	} else {
		path = strings.Replace(path, "{owner}", owner, 1)
		path = strings.Replace(path, "{repo}", repo, 1)
	}

	gh := github.NewClient(host)
	response, err := gh.GenericAPIRequest(method, path, params, cacheTTL)
	utils.Check(err)
	defer response.Body.Close()

	colorize := ui.IsTerminal(os.Stdout)
	if args.Flag.Bool("--flat") {
		utils.JSONPath(ui.Stdout, response.Body, colorize)
	} else {
		io.Copy(ui.Stdout, response.Body)
		ui.Println()
	}

	args.NoForward()
}

func valueOrFileContents(value string) string {
	if strings.HasPrefix(value, "@") {
		file := strings.TrimPrefix(value, "@")
		var content []byte
		var err error
		if file == "-" {
			content, err = ioutil.ReadAll(os.Stdin)
		} else {
			content, err = ioutil.ReadFile(file)
		}
		if err != nil {
			utils.Check(err)
		}
		return string(content)
	} else {
		return value
	}
}

func quote(s string) string {
	return fmt.Sprintf("%q", s)
}
