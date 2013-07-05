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
	for i, a := range args.Params {
		if hasValueRegxp.MatchString(a) {
			continue
		}

		if nameWithOwnerRegexp.MatchString(a) && !isDir(a) {
			name, owner := parseCloneNameAndOwner(a)
			config := github.CurrentConfig()
			isSSH = isSSH || owner == config.User
			if owner == "" {
				owner = config.User
			}

			project := github.Project{Name: name, Owner: owner}
			url := project.GitURL(name, owner, isSSH)
			args.ReplaceParam(i, url)

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
