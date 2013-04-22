package main

import (
	"bytes"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

var git = Git{}

type Git struct{}

func (g *Git) Dir() string {
	return execGitCmd("rev-parse -q --git-dir")[0]
}

func (g *Git) Editor() string {
	return execGitCmd("var GIT_EDITOR")[0]
}

func (g *Git) Owner() string {
	remote := g.Remote()
	return mustMatchGitUrl(remote)[1]
}

func (g *Git) Repo() string {
	remote := g.Remote()
	return mustMatchGitUrl(remote)[2]
}

func (g *Git) CurrentBranch() string {
	return execGitCmd("symbolic-ref -q --short HEAD")[0]
}

// FIXME: only care about origin push remote now
func (g *Git) Remote() string {
	r := regexp.MustCompile("origin\t(.+) \\(push\\)")
	for _, output := range execGitCmd("remote -v") {
		if r.MatchString(output) {
			return r.FindStringSubmatch(output)[1]
		}
	}

	panic("Can't find remote")
}

func execGitCmd(args string) []string {
	return execCmd("git " + args)
}

func execCmd(input string) (outputs []string) {
	inputs := strings.Split(input, " ")
	name := inputs[0]
	args := inputs[1:]

	out, err := exec.Command(name, args...).Output()
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range bytes.Split(out, []byte{'\n'}) {
		outputs = append(outputs, string(line))
	}

	return outputs
}
