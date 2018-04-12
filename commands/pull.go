package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var cmdPull = &Command{
	Run:          pull,
	GitExtension: true,
	Usage:        "pull [-p] [<OPTIONS>] [<USER>/]<REPOSITORY> [<DESTINATION>]",
	Long: `Pull a repository from GitHub.

## Options:
	-p
		(Deprecated) Pull private repositories over SSH.

	[<USER>/]<REPOSITORY>
		<USER> defaults to your own GitHub username.

	<DESTINATION>
		Directory name to pull into (default: <REPOSITORY>).

## Examples:
		$ hub pull rtomayko/ronn
		> git pull git://github.com/rtomayko/ronn.git

## See also:

hub-fork(1), hub(1), git-pull(1)
`,
}

func init() {
	CmdRunner.Use(cmdPull)
}

func pull(command *Command, args *Args) {
	if !args.IsParamsEmpty() {
		transformPullArgs(args)
	}
}

func transformPullArgs(args *Args) {
	isSSH := parsePullPrivateFlag(args)
	hasValueRegexp := regexp.MustCompile("^(--(upload-pack|template|depth|origin|branch|reference|name)|-[ubo])$")
	nameWithOwnerRegexp := regexp.MustCompile(NameWithOwnerRe)
	for i := 0; i < args.ParamsSize(); i++ {
		a := args.Params[i]

		if strings.HasPrefix(a, "-") {
			if hasValueRegexp.MatchString(a) {
				i++
			}
		} else {
			if nameWithOwnerRegexp.MatchString(a) && !isCloneable(a) {
				name, owner := parsePullNameAndOwner(a)
				var host *github.Host
				if owner == "" {
					config := github.CurrentConfig()
					h, err := config.DefaultHost()
					if err != nil {
						utils.Check(github.FormatError("pulling repository", err))
					}

					host = h
					owner = host.User
				}

				var hostStr string
				if host != nil {
					hostStr = host.Host
				}

				expectWiki := strings.HasSuffix(name, ".wiki")
				if expectWiki {
					name = strings.TrimSuffix(name, ".wiki")
				}

				project := github.NewProject(owner, name, hostStr)
				gh := github.NewClient(project.Host)
				repo, err := gh.Repository(project)
				if err != nil {
					if strings.Contains(err.Error(), "HTTP 404") {
						err = fmt.Errorf("Error: repository %s/%s doesn't exist", project.Owner, project.Name)
					}
					utils.Check(err)
				}

				owner = repo.Owner.Login
				name = repo.Name
				if expectWiki {
					if !repo.HasWiki {
						utils.Check(fmt.Errorf("Error: %s/%s doesn't have a wiki", owner, name))
					} else {
						name = name + ".wiki"
					}
				}

				if !isSSH &&
					args.Command != "submodule" &&
					!github.IsHttpsProtocol() {
					isSSH = repo.Private || repo.Permissions.Push
				}

				url := project.GitURL(name, owner, isSSH)
				args.ReplaceParam(i, url)
			}

			break
		}
	}
}

func parsePullPrivateFlag(args *Args) bool {
	if i := args.IndexOfParam("-p"); i != -1 {
		args.RemoveParam(i)
		return true
	}

	return false
}

func parsePullNameAndOwner(arg string) (name, owner string) {
	name, owner = arg, ""
	if strings.Contains(arg, "/") {
		split := strings.SplitN(arg, "/", 2)
		name = split[1]
		owner = split[0]
	}

	return
}
