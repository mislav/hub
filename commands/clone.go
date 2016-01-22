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
	Usage:        "clone [-p] OPTIONS [USER/]REPOSITORY [DIRECTORY]",
	Short:        "Clone a remote repository into a new directory",
	Long: `Clone repository "git://github.com/USER/REPOSITORY.git" into
DIRECTORY as with git-clone(1). When USER/ is omitted, assumes
your GitHub login. With -p, clone private repositories over SSH.
For repositories under your GitHub login, -p is implicit.
`,
}

func init() {
	CmdRunner.Use(cmdClone)
}

/**
  $ hub clone jingweno/gh
  > git clone git://github.com/jingweno/gh.git

  $ hub clone -p jingweno/gh
  > git clone git@github.com:jingweno/gh.git

  $ hub clone jekyll_and_hyde
  > git clone git://github.com/YOUR_LOGIN/jekyll_and_hyde.git

  $ hub clone -p jekyll_and_hyde
  > git clone git@github.com:YOUR_LOGIN/jekyll_and_hyde.git

  $ hub clone jekyll_and_hyde jekyll
  > git clone git://github.com/YOUR_LOGIN/jekyll_and_hyde.git jekyll
*/
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

				if expectWiki {
					if !repo.HasWiki {
						utils.Check(fmt.Errorf("Error: %s/%s doesn't have a wiki", project.Owner, project.Name))
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
