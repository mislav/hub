package commands

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/github/hub/v2/github"
	"github.com/github/hub/v2/ui"
	"github.com/github/hub/v2/utils"
)

var cmdAPI = &Command{
	Run:   apiCommand,
	Usage: "api [-it] [-X <METHOD>] [-H <HEADER>] [--cache <TTL>] <ENDPOINT> [-F <FIELD>|--input <FILE>]",
	Long: `Low-level GitHub API request interface.

## Options:
	-X, --method <METHOD>
		The HTTP method to use for the request (default: "GET"). The method is
		automatically set to "POST" if ''--field'', ''--raw-field'', or ''--input''
		are used.

		Use ''-XGET'' to force serializing fields into the query string for the GET
		request instead of JSON body of the POST request.

	-F, --field <KEY>=<VALUE>
		Data to serialize with the request. <VALUE> has some magic handling; use
		''--raw-field'' for sending arbitrary string values.

		If <VALUE> starts with "@", the rest of the value is interpreted as a
		filename to read the value from. Use "@-" to read from standard input.

		If <VALUE> is "true", "false", "null", or looks like a number, an
		appropriate JSON type is used instead of a string.

		It is not possible to serialize <VALUE> as a nested JSON array or hash.
		Instead, construct the request payload externally and pass it via
		''--input''.

		Unless ''-XGET'' was used, all fields are sent serialized as JSON within
		the request body. When <ENDPOINT> is "graphql", all fields other than
		"query" are grouped under "variables". See
		<https://graphql.org/learn/queries/#variables>

	-f, --raw-field <KEY>=<VALUE>
		Same as ''--field'', except that it allows values starting with "@", literal
		strings "true", "false", and "null", as well as strings that look like
		numbers.

	--input <FILE>
		The filename to read the raw request body from. Use "-" to read from standard
		input. Use this when you want to manually construct the request payload.

	-H, --header <KEY>:<VALUE>
		Set an HTTP request header.

	-i, --include
		Include HTTP response headers in the output.

	-t, --flat
		Parse response JSON and output the data in a line-based key-value format
		suitable for use in shell scripts.

	--paginate
		Automatically request and output the next page of results until all
		resources have been listed. For GET requests, this follows the ''<next\>''
		resource as indicated in the "Link" response header. For GraphQL queries,
		this utilizes ''pageInfo'' that must be present in the query; see EXAMPLES.

		Note that multiple JSON documents will be output as a result. If the API
		rate limit has been reached, the final document that is output will be the
		HTTP 403 notice, and the process will exit with a non-zero status. One way
		this can be avoided is by enabling ''--obey-ratelimit''.

	--color[=<WHEN>]
		Enable colored output even if stdout is not a terminal. <WHEN> can be one
		of "always" (default for ''--color''), "never", or "auto" (default).

	--cache <TTL>
		Cache valid responses to GET requests for <TTL> seconds.

		When using "graphql" as <ENDPOINT>, caching will apply to responses to POST
		requests as well. Just make sure to not use ''--cache'' for any GraphQL
		mutations.

	--obey-ratelimit
		After exceeding the API rate limit, pause the process until the reset time
		of the current rate limit window and retry the request. Note that this may
		cause the process to hang for a long time (maximum of 1 hour).

	<ENDPOINT>
		The GitHub API endpoint to send the HTTP request to (default: "/").

		To learn about available endpoints, see <https://developer.github.com/v3/>.
		To make GraphQL queries, use "graphql" as <ENDPOINT> and pass ''-F query=QUERY''.

		If the literal strings "{owner}" or "{repo}" appear in <ENDPOINT> or in the
		GraphQL "query" field, fill in those placeholders with values read from the
		git remote configuration of the current git repository.

## Examples:

		# fetch information about the currently authenticated user as JSON
		$ hub api user

		# list user repositories as line-based output
		$ hub api --flat users/octocat/repos

		# post a comment to issue #23 of the current repository
		$ hub api repos/{owner}/{repo}/issues/23/comments --raw-field 'body=Nice job!'

		# perform a GraphQL query read from a file
		$ hub api graphql -F query=@path/to/myquery.graphql

		# perform pagination with GraphQL
		$ hub api --paginate graphql -f query='
		  query($endCursor: String) {
		    repositoryOwner(login: "USER") {
		      repositories(first: 100, after: $endCursor) {
		        nodes {
		          nameWithOwner
		        }
		        pageInfo {
		          hasNextPage
		          endCursor
		        }
		      }
		    }
		  }
		'

## See also:

hub(1)
`,
}

func init() {
	CmdRunner.Use(cmdAPI)
}

func apiCommand(_ *Command, args *Args) {
	path := ""
	if !args.IsParamsEmpty() {
		path = args.GetParam(0)
	}

	method := "GET"
	if args.Flag.HasReceived("--method") {
		method = args.Flag.Value("--method")
	} else if args.Flag.HasReceived("--field") || args.Flag.HasReceived("--raw-field") || args.Flag.HasReceived("--input") {
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

	isGraphQL := path == "graphql"
	if isGraphQL && params["query"] != nil {
		query := params["query"].(string)
		query = strings.Replace(query, "{owner}", owner, -1)
		query = strings.Replace(query, "{repo}", repo, -1)

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
		path = strings.Replace(path, "{owner}", owner, -1)
		path = strings.Replace(path, "{repo}", repo, -1)
	}

	var body interface{}
	if args.Flag.HasReceived("--input") {
		fn := args.Flag.Value("--input")
		if fn == "-" {
			body = os.Stdin
		} else {
			fi, err := os.Open(fn)
			utils.Check(err)
			body = fi
			defer fi.Close()
		}
	} else {
		body = params
	}

	gh := github.NewClient(host)

	out := ui.Stdout
	colorize := colorizeOutput(args.Flag.HasReceived("--color"), args.Flag.Value("--color"))
	parseJSON := args.Flag.Bool("--flat")
	includeHeaders := args.Flag.Bool("--include")
	paginate := args.Flag.Bool("--paginate")
	rateLimitWait := args.Flag.Bool("--obey-ratelimit")

	args.NoForward()

	for {
		response, err := gh.GenericAPIRequest(method, path, body, headers, cacheTTL)
		utils.Check(err)

		if rateLimitWait && response.StatusCode == 403 && response.RateLimitRemaining() == 0 {
			pauseUntil(response.RateLimitReset())
			continue
		}

		success := response.StatusCode < 300
		jsonType := true
		if !success {
			jsonType, _ = regexp.MatchString(`[/+]json(?:;|$)`, response.Header.Get("Content-Type"))
		}

		if includeHeaders {
			fmt.Fprintf(out, "%s %s\r\n", response.Proto, response.Status)
			response.Header.Write(out)
			fmt.Fprintf(out, "\r\n")
		}

		endCursor := ""
		hasNextPage := false

		if parseJSON && jsonType {
			hasNextPage, endCursor = utils.JSONPath(out, response.Body, colorize)
		} else if paginate && isGraphQL {
			bodyCopy := &bytes.Buffer{}
			io.Copy(out, io.TeeReader(response.Body, bodyCopy))
			hasNextPage, endCursor = utils.JSONPath(ioutil.Discard, bodyCopy, false)
		} else {
			io.Copy(out, response.Body)
		}
		response.Body.Close()

		if !success {
			if ssoErr := github.ValidateGitHubSSO(response.Response); ssoErr != nil {
				ui.Errorln()
				ui.Errorln(ssoErr)
			}
			if scopeErr := github.ValidateSufficientOAuthScopes(response.Response); scopeErr != nil {
				ui.Errorln()
				ui.Errorln(scopeErr)
			}
			os.Exit(22)
		}

		if paginate {
			if isGraphQL && hasNextPage && endCursor != "" {
				if v, ok := params["variables"]; ok {
					variables := v.(map[string]interface{})
					variables["endCursor"] = endCursor
				} else {
					variables := map[string]interface{}{"endCursor": endCursor}
					params["variables"] = variables
				}
				goto next
			} else if nextLink := response.Link("next"); nextLink != "" {
				path = nextLink
				goto next
			}
		}

		break
	next:
		if !parseJSON {
			fmt.Fprintf(out, "\n")
		}

		if rateLimitWait && response.RateLimitRemaining() == 0 {
			pauseUntil(response.RateLimitReset())
		}
	}
}

func pauseUntil(timestamp int) {
	rollover := time.Unix(int64(timestamp)+1, 0)
	duration := time.Until(rollover)
	if duration > 0 {
		ui.Errorf("API rate limit exceeded; pausing until %v ...\n", rollover)
		time.Sleep(duration)
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
