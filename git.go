package main

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

func FetchGitDir() string {
	gitDir := execGitCmd([]string{"rev-parse", "-q", "--git-dir"})[0]
	gitDir, err := filepath.Abs(gitDir)
	if err != nil {
		log.Fatal(err)
	}

	return gitDir
}

func FetchGitEditor() string {
	return execGitCmd([]string{"var", "GIT_EDITOR"})[0]
}

func FetchGitOwner() string {
	remote := FetchGitRemote()
	return mustMatchGitUrl(remote)[1]
}

func FetchGitProject() string {
	remote := FetchGitRemote()
	return mustMatchGitUrl(remote)[2]
}

func FetchGitHead() string {
	return execGitCmd([]string{"symbolic-ref", "-q", "--short", "HEAD"})[0]
}

// FIXME: only care about origin push remote now
func FetchGitRemote() string {
	r := regexp.MustCompile("origin\t(.+) \\(push\\)")
	for _, output := range execGitCmd([]string{"remote", "-v"}) {
		if r.MatchString(output) {
			return r.FindStringSubmatch(output)[1]
		}
	}

	panic("Can't find remote")
}

func FetchGitCommitLogs(sha1, sha2 string) string {
	execCmd := NewExecCmd("git")
	execCmd.WithArg("log").WithArg("--no-color")
	execCmd.WithArg("--format=%h (%aN, %ar)%n%w(78,3,3)%s%n%+b")
	execCmd.WithArg("--cherry")
	shaRange := fmt.Sprintf("%s...%s", sha1, sha2)
	execCmd.WithArg(shaRange)

	outputs, err := execCmd.ExecOutput()
	if err != nil {
		log.Fatal(err)
	}

	return outputs
}

func execGitCmd(input []string) (outputs []string) {
	cmd := NewExecCmd("git")
	for _, i := range input {
		cmd.WithArg(i)
	}

	out, err := cmd.ExecOutput()
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range strings.Split(out, "\n") {
		outputs = append(outputs, string(line))
	}

	return outputs
}
