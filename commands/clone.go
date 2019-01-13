package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/github/hub/cmd"
	"github.com/github/hub/github"
	"github.com/github/hub/ui"
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

## Protocol used for cloning

The 'git:' protocol will be used for cloning public repositories, while the SSH
protocol will be used for private repositories and those that you have push
access to. Alternatively, hub can be configured to use HTTPS protocol for
everything. See "HTTPS instead of git protocol" and "HUB_PROTOCOL" of hub(1).

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
	hasValueRegexp := regexp.MustCompile("^(--(upload-pack|template|depth|origin|branch|reference|name)|-[ubo])$")
	nameWithOwnerRegexp := regexp.MustCompile(NameWithOwnerRe)
	var targetDir string
	var originName string
	quiet := false
	for i := 0; i < args.ParamsSize(); i++ {
		a := args.Params[i]

		if strings.HasPrefix(a, "-") {
			if a == "--origin" || a == "-o" {
				if i+1 < args.ParamsSize() {
					originName = args.Params[i+1]
				}
			} else if strings.HasPrefix(a, "--origin=") {
				originName = strings.TrimPrefix(a, "--origin=")
			} else if strings.HasPrefix(a, "-o") {
				originName = strings.TrimPrefix(a, "-o")
			} else if a == "--quiet" || a == "-q" {
				quiet = true
			}
			if hasValueRegexp.MatchString(a) {
				i++
			}
		} else {
			if targetDir != "" {
				targetDir = a
				break
			} else if nameWithOwnerRegexp.MatchString(a) && !isCloneable(a) {
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

				targetDir = name
				if repo.Parent != nil && args.Command != "submodule" {
					args.AfterFn(func() error {
						upstreamName := "upstream"
						if originName == "upstream" {
							upstreamName = "origin"
						}
						upstreamUrl := project.GitURL(repo.Parent.Name, repo.Parent.Owner.Login, repo.Parent.Private)
						addRemote := cmd.New("git")
						addRemote.WithArgs("-C", targetDir)
						addRemote.WithArgs("remote", "add", "-f", upstreamName, upstreamUrl)
						if !quiet {
							ui.Errorf("Adding remote '%s' for the '%s/%s' repo\n",
								upstreamName, repo.Parent.Owner.Login, repo.Parent.Name)
						}
						output, err := addRemote.CombinedOutput()
						if err != nil {
							ui.Errorln(output)
						}
						return err
					})
				}
			} else {
				break
			}
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
