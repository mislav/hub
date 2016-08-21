package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var cmdClone = &Command{
	Run:          clone,
	GitExtension: true,
	Usage:        "clone [-p] [<OPTIONS>] [<USER>/]<REPOSITORY> [<DESTINATION>]",
	Long: `Clone a repository from GitHub.

## Options:
	-p
		(Deprecated) Clone private repositories over SSH.

	[<USER>/]<REPOSITORY>
		<USER> defaults to your own GitHub username.

	<DESTINATION>
		Directory name to clone into (default: <REPOSITORY>).

## Examples:
		$ hub clone rtomayko/ronn
		> git clone git://github.com/rtomayko/ronn.git

## See also:

hub-fork(1), hub(1), git-clone(1)
`,
}

func init() {
	CmdRunner.Use(cmdClone)
}

func clone(command *Command, args *Args) {
	if !args.IsParamsEmpty() {
		transformCloneArgs(args)
	}
}

func transformCloneArgs(args *Args) {
	isSSH := parseClonePrivateFlag(args)
	hasValueRegxp := regexp.MustCompile("^(--(upload-pack|template|depth|origin|branch|reference|name)|-[ubo])$")
	nameWithOwnerRegexp := regexp.MustCompile(NameWithOwnerRe)
	for i := 0; i < args.ParamsSize(); i++ {
		a := args.Params[i]

		if strings.HasPrefix(a, "-") {
			if hasValueRegxp.MatchString(a) {
				i++
			}
		} else {
			if nameWithOwnerRegexp.MatchString(a) && !isCloneable(a) {
				name, owner := parseCloneNameAndOwner(a)
				var host *github.Host
				if owner == "" {
					config := github.CurrentConfig()
					h, err := config.DefaultHost()
					if err != nil {
						utils.Check(github.FormatError("cloning repository", err))
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

func parseClonePrivateFlag(args *Args) bool {
	if i := args.IndexOfParam("-p"); i != -1 {
		args.RemoveParam(i)
		return true
	}

	return false
}

func parseCloneNameAndOwner(arg string) (name, owner string) {
	name, owner = arg, ""
	if strings.Contains(arg, "/") {
		split := strings.SplitN(arg, "/", 2)
		name = split[1]
		owner = split[0]
	}

	return
}
