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

	// git help clone | grep -e '^ \+-.\+<'
	p := utils.NewArgsParser()
	p.RegisterValue("--branch", "-b")
	p.RegisterValue("--depth")
	p.RegisterValue("--reference")
	if args.Command == "submodule" {
		p.RegisterValue("--name")
	} else {
		p.RegisterValue("--config", "-c")
		p.RegisterValue("--jobs", "-j")
		p.RegisterValue("--origin", "-o")
		p.RegisterValue("--reference-if-able")
		p.RegisterValue("--separate-git-dir")
		p.RegisterValue("--shallow-exclude")
		p.RegisterValue("--shallow-since")
		p.RegisterValue("--template")
		p.RegisterValue("--upload-pack", "-u")
		p.RegisterBool("--quiet", "-q")
	}
	p.Parse(args.Params)

	upstreamName := "upstream"
	originName := p.Value("--origin")
	quiet := p.Bool("--quiet")
	targetDir := ""

	nameWithOwnerRegexp := regexp.MustCompile(NameWithOwnerRe)
	for n, i := range p.PositionalIndices {
		switch n {
		case 0:
			repo := args.Params[i]
			if nameWithOwnerRegexp.MatchString(repo) && !isCloneable(repo) {
				name := repo
				owner := ""
				if strings.Contains(name, "/") {
					split := strings.SplitN(name, "/", 2)
					owner = split[0]
					name = split[1]
				}

				config := github.CurrentConfig()
				host, err := config.DefaultHost()
				if err != nil {
					utils.Check(github.FormatError("cloning repository", err))
				}
				if owner == "" {
					owner = host.User
				}

				expectWiki := strings.HasSuffix(name, ".wiki")
				if expectWiki {
					name = strings.TrimSuffix(name, ".wiki")
				}

				project := github.NewProject(owner, name, host.Host)
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

				targetDir = name
				url := project.GitURL(name, owner, isSSH)
				args.ReplaceParam(i, url)

				if repo.Parent != nil && args.Command == "clone" && originName != upstreamName {
					args.AfterFn(func() error {
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
				} else {
					break
				}
			}
		case 1:
			targetDir = args.Params[i]
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
