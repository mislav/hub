package commands

import (
	"github.com/jingweno/gh/github"
	"regexp"
	"strings"
	"os"
	"fmt"
)

var cmdSubmodule = &Command{
	Run: submodule,
	GitExtension: true,
	Usage: "submodule add [-p] OPTIONS [USER/]REPOSITORY DIRECTORY",
	Short: "Initialize, update or inspect submodules",
	Long: `Submodule repository "git://github.com/USER/REPOSITORY.git" into
DIRECTORY as  with  git-submodule(1).  When  USER/  is  omitted,
assumes   your   GitHub  login.  With  -p,  use  private  remote
"git@github.com:USER/REPOSITORY.git".`,
}

/**
 $ gh submodule add jingweno/gh vendor/gh
 > git submodule add git://github.com/jingweno/gh.git vendor/gh

 $ gh submodule add -p jingweno/gh vendor/gh
 > git submodule add git@github.com:jingweno/gh.git vendor/gh

 $ gh submodule add -b gh --name gh jingweno/gh vendor/gh
 > git submodule add -b gh --name gh git://github.com/jingweno/gh.git vendor/gh
**/

func submodule(command *Command, args *Args) {
	if !args.IsParamsEmpty() {
		transformSubmoduleArgs(args)
	}
}

func transformSubmoduleArgs(args *Args) {
	isSSH := parseSubmodulePrivateFlag(args)
	
	nameWithOwnerRegexp := regexp.MustCompile(NameWithOwnerRe)
	hasValueRegexp := regexp.MustCompile("^(--(reference|name)|-b)$")
	
	var continueNext bool

	for i, a := range args.Params {
		if continueNext {
			continueNext = false
			continue
		}

		if hasValueRegexp.MatchString(a) {
			if !strings.Contains(a, "=") {
				continueNext = true
			}

			continue
		}

		if nameWithOwnerRegexp.MatchString(a) && !isDir(a) && a != "add" {
			name, owner := parseSubmoduleNameAndOwner(a)
			config := github.CurrentConfig()
			isSSH = isSSH || owner == config.User
			if owner == "" {
				owner = config.User
			}

			project := github.Project{Name: name, Owner: owner}
			url := project.GitURL(name, owner, isSSH)
			
			args.ReplaceParam(i, url)

			if args.Noop {
				fmt.Printf("it would run `git submodule %s`\n", strings.Join(args.Params, " "))
				os.Exit(0)
			}

			break
		}
	}
}

func parseSubmodulePrivateFlag(args *Args) bool {
	if i := args.IndexOfParam("-p"); i != -1 {
		args.RemoveParam(i)
		return true
	}

	return false
}

func parseSubmoduleNameAndOwner(arg string) (name, owner string) {
	name, owner = arg, ""
	if strings.Contains(arg, "/") {
		split := strings.SplitN(arg, "/", 2)
		name = split[1]
		owner = split[0]
	}

	return
}
