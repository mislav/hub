package commands

import (
	"github.com/jingweno/gh/github"
	"regexp"
	"strings"
)

var cmdClone = &Command{
	Run:          clone,
	GitExtension: true,
	Usage:        "clone [-p] OPTIONS [USER/]REPOSITORY DIRECTORY",
	Short:        "Clone a remote repository into a new directory",
	Long: `Clone repository "git://github.com/USER/REPOSITORY.git" into
DIRECTORY as with git-clone(1). When USER/ is omitted, assumes
your GitHub login. With -p, clone private repositories over SSH.
For repositories under your GitHub login, -p is implicit.
`,
}

/**
  $ gh clone jingweno/gh
  > git clone git://github.com/jingweno/gh.git

  $ gh clone -p jingweno/gh
  > git clone git@github.com:jingweno/gh.git

  $ gh clone jekyll_and_hyde
  > git clone git://github.com/YOUR_LOGIN/jekyll_and_hyde.git

  $ gh clone -p jekyll_and_hyde
  > git clone git@github.com:YOUR_LOGIN/jekyll_and_hyde.git
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
			if nameWithOwnerRegexp.MatchString(a) && !isDir(a) {
				name, owner := parseCloneNameAndOwner(a)
				var credentials *github.Credentials
				if owner == "" {
					configs := github.CurrentConfigs()
					credentials = configs.DefaultCredentials()
					owner = credentials.User
				}

				var host string
				if credentials != nil {
					host = credentials.Host
				}

				project := github.NewProject(owner, name, host)
				isSSH = isSSH ||
					args.Command != "submodule" &&
						credentials != nil &&
						project.Owner == credentials.User

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
