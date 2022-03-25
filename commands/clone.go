package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/github/hub/v2/github"
	"github.com/github/hub/v2/utils"
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

HTTPS protocol is used by hub as the default. Alternatively, hub can be
configured to use SSH protocol for all git operations. See "SSH instead
of HTTPS protocol" and "HUB_PROTOCOL" of hub(1).

## Examples:
		$ hub clone rtomayko/ronn
		> git clone https://github.com/rtomayko/ronn.git

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
	isPrivate := parseClonePrivateFlag(args)

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
	}
	p.Parse(args.Params)

	nameWithOwnerRegexp := regexp.MustCompile(NameWithOwnerRe)
	if len(p.PositionalIndices) > 0 {
		i := p.PositionalIndices[0]
		a := args.Params[i]
		if nameWithOwnerRegexp.MatchString(a) && !isCloneable(a) {
			url := getCloneURL(a, isPrivate, args.Command != "submodule")
			args.ReplaceParam(i, url)
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

func getCloneURL(nameWithOwner string, allowPush, allowPrivate bool) string {
	name := nameWithOwner
	owner := ""
	if strings.Contains(name, "/") {
		split := strings.SplitN(name, "/", 2)
		owner = split[0]
		name = split[1]
	}

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

	if !allowPush && allowPrivate {
		allowPush = repo.Private || repo.Permissions.Push
	}

	return project.GitURL(name, owner, allowPush)
}
