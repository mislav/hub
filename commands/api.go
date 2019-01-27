package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

var cmdApi = &Command{
	Run:   apiCommand,
	Usage: "api [-t] [-X <METHOD>] [--cache <TTL>] <ENDPOINT> [-F <KEY-VALUE>]",
	Long: `Low-level GitHub API request interface.

## Options:
	-X, --method <METHOD>
		The HTTP method to use for the request (default: "GET"). The method is
		automatically set to "POST" if '--field' or '--raw-field' are used.

		Use '-XGET' to force serializing fields into the query string for the GET
		request instead of JSON body of the POST request.

	-F, --field <KEY-VALUE>
		Send data in 'KEY=VALUE' format. The <VALUE> part has some magic handling;
		see '--raw-field' for passing arbitrary strings.

		If <VALUE> starts with "@", the rest of the value is interpreted as a
		filename to read the value from. Use "@-" to read from standard input.

		If <VALUE> is "true", "false", "null", or looks like a number, an
		appropriate JSON type is used instead of a string.

		Unless '-XGET' was used, all fields are sent serialized as JSON within the
		request body. When <ENDPOINT> is "graphql", all fields other than "query"
		are grouped under "variables". See
		<https://graphql.org/learn/queries/#variables>

	-f, --raw-field <KEY-VALUE>
		Same as '--field', except that it allows values starting with "@", literal
		strings "true", "false", and "null", as well as strings that look like
		numbers.

	-H, --header <KEY-VALUE>
		An HTTP request header in 'KEY: VALUE' format.

	-i, --include
		Include HTTP response headers in the output.

	-t, --flat
		Parse response JSON and output the data in a line-based key-value format
		suitable for use in shell scripts.

	--cache <TTL>
		Cache successful responses to GET requests for <TTL> seconds.

		When using "graphql" as <ENDPOINT>, caching will apply to responses to POST
		requests as well. Just make sure to not use '--cache' for any GraphQL
		mutations.

	<ENDPOINT>
		The GitHub API endpoint to send the HTTP request to (default: "/").
		
		To learn about available endpoints, see <https://developer.github.com/v3/>.
		To make GraphQL queries, use "graphql" as <ENDPOINT> and pass '-F query=QUERY'.

		If the literal strings "{owner}" or "{repo}" appear in <ENDPOINT> or in the
		GraphQL "query" field, fill in those placeholders with values read from the
		git remote configuration of the current git repository.

## Examples:

		# fetch information about the currently authenticated user as JSON
		$ hub api user

		# list user repositories as line-based output
		$ hub api --flat users/octocat/repos

		# post a comment to issue #23 of the current repository
		$ hub api repos/{owner}/{repo}/issues/23/comments --raw-field "body=Nice job!"

		# perform a GraphQL query read from a file
		$ hub api graphql -F query=@path/to/myquery.graphql

## See also:

hub(1)
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
	} else if args.Flag.HasReceived("--field") || args.Flag.HasReceived("--raw-field") {
		method = "POST"
	}
	cacheTTL := args.Flag.Int("--cache")

	params := make(map[string]interface{})
	for _, val := range args.Flag.AllValues("--field") {
		parts := strings.SplitN(val, "=", 2)
		if len(parts) >= 2 {
			params[parts[0]] = magicValue(parts[1])
		}
	}
	for _, val := range args.Flag.AllValues("--raw-field") {
		parts := strings.SplitN(val, "=", 2)
		if len(parts) >= 2 {
			params[parts[0]] = parts[1]
		}
	}

	headers := make(map[string]string)
	for _, val := range args.Flag.AllValues("--header") {
		parts := strings.SplitN(val, ":", 2)
		if len(parts) >= 2 {
			headers[parts[0]] = strings.TrimLeft(parts[1], " ")
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
		defHost, err := github.CurrentConfig().DefaultHostNoPrompt()
		utils.Check(err)
		host = defHost.Host
	}

	if path == "graphql" && params["query"] != nil {
		query := params["query"].(string)
		query = strings.Replace(query, quote("{owner}"), quote(owner), 1)
		query = strings.Replace(query, quote("{repo}"), quote(repo), 1)

		variables := make(map[string]interface{})
		for key, value := range params {
			if key != "query" {
				variables[key] = value
			}
		}
		if len(variables) > 0 {
			params = make(map[string]interface{})
			params["variables"] = variables
		}

		params["query"] = query
	} else {
		path = strings.Replace(path, "{owner}", owner, 1)
		path = strings.Replace(path, "{repo}", repo, 1)
	}

	gh := github.NewClient(host)
	response, err := gh.GenericAPIRequest(method, path, params, headers, cacheTTL)
	utils.Check(err)
	defer response.Body.Close()

	args.NoForward()

	out := ui.Stdout
	colorize := ui.IsTerminal(os.Stdout)
	success := response.StatusCode < 300
	parseJSON := args.Flag.Bool("--flat")

	if !success {
		jsonType, _ := regexp.MatchString(`[/+]json(?:;|$)`, response.Header.Get("Content-Type"))
		parseJSON = parseJSON && jsonType
	}

	if args.Flag.Bool("--include") {
		fmt.Fprintf(out, "%s %s\r\n", response.Proto, response.Status)
		response.Header.Write(out)
		fmt.Fprintf(out, "\r\n")
	}

	if parseJSON {
		utils.JSONPath(out, response.Body, colorize)
	} else {
		io.Copy(out, response.Body)
	}

	if !success {
		os.Exit(22)
	}
}

const (
	trueVal  = "true"
	falseVal = "false"
	nilVal   = "null"
)

func magicValue(value string) interface{} {
	switch value {
	case trueVal:
		return true
	case falseVal:
		return false
	case nilVal:
		return nil
	default:
		if strings.HasPrefix(value, "@") {
			return string(readFile(value[1:]))
		} else if i, err := strconv.Atoi(value); err == nil {
			return i
		} else {
			return value
		}
	}
}

func readFile(file string) (content []byte) {
	var err error
	if file == "-" {
		content, err = ioutil.ReadAll(os.Stdin)
	} else {
		content, err = ioutil.ReadFile(file)
	}
	utils.Check(err)
	return
}

func quote(s string) string {
	return fmt.Sprintf("%q", s)
}
