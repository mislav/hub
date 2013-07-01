package commands

import (
	"github.com/jingweno/gh/github"
	"regexp"
	"strings"
)

var cmdClone = &Command{
	Run:          clone,
	GitExtension: true,
	Usage:        "clone [-p] OPTIONS [USER] REPOSITORY DIRECTORY",
	Short:        "clone a remote repository into a new directory",
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
			owner, name := "", a
			if strings.Contains(a, "/") {
				split := strings.SplitN(a, "/", 2)
				owner = split[0]
				name = split[1]
			}

			if owner == "" {
				config := github.CurrentConfig()
				owner = config.User
			}

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
