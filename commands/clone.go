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
	Short:        "clone a remote repository into a new directory",
	Long: `Clone repository "git://github.com/USER/REPOSITORY.git" into
DIRECTORY as with git-clone(1). When USER/ is omitted, assumes
your GitHub login. With -p, clone private repositories over SSH.
For repositories under your GitHub login, -p is implicit.
`,
}

/**
  $ gh clone jingweno/gh
  > git clone git://github.com/jingweno/gh

  $ gh clone -p jingweno/gh
  > git clone git@github.com:jingweno/gh.git

  $ gh clone jekyll_and_hype
  > git clone git://github.com/YOUR_LOGIN/jekyll_and_hype.

  $ hub clone -p jekyll_and_hype
  > git clone git@github.com:YOUR_LOGIN/jekyll_and_hype.git
*/
func clone(command *Command, args *Args) {
	if !args.IsEmpty() {
		transformCloneArgs(args)
	}
}

func transformCloneArgs(args *Args) {
	isSSH := parseClonePrivateFlag(args)
	hasValueRegxp := regexp.MustCompile("^(--(upload-pack|template|depth|origin|branch|reference|name)|-[ubo])$")
	nameWithOwnerRegexp := regexp.MustCompile(NameWithOwnerRe)
	for i, a := range args.Array() {
		if hasValueRegxp.MatchString(a) {
			continue
		}

		if nameWithOwnerRegexp.MatchString(a) && !isDir(a) {
			name, owner := parseCloneNameAndOwner(a)
			config := github.CurrentConfig()
			if owner == "" {
				owner = config.User
			}
			isSSH = isSSH || owner == config.User

			project := github.Project{Name: name, Owner: owner}
			url := project.GitURL(name, owner, isSSH)
			args.Replace(i, url)

			break
		}
	}
}

func parseClonePrivateFlag(args *Args) bool {
	if i := args.IndexOf("-p"); i != -1 {
		args.Remove(i)
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
