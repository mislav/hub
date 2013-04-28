package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var git = Git{}

type Git struct{}

func (g *Git) Dir() string {
	gitDir := execGitCmd("rev-parse -q --git-dir")[0]
	gitDir, err := filepath.Abs(gitDir)
	if err != nil {
		log.Fatal(err)
	}

	return gitDir
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

func (g *Git) CommitLogs(sha1, sha2 string) string {
	execCmd := NewExecCmd("git")
	execCmd.WithArg("log").WithArg("--no-color")
	execCmd.WithArg("--format=%h (%aN, %ar)%n%w(78,3,3)%s%n%+b")
	execCmd.WithArg("--cherry")
	shaRange := fmt.Sprintf("%s...%s", sha1, sha2)
	execCmd.WithArg(shaRange)

	outputs, err := execCmd.Exec()
	if err != nil {
		log.Fatal(err)
	}

	return outputs
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

func execGitCmd(input string) (outputs []string) {
	name := "git"
	args := strings.Split(input, " ")

	out, err := exec.Command(name, args...).Output()
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range bytes.Split(out, []byte{'\n'}) {
		outputs = append(outputs, string(line))
	}

	return outputs
}
