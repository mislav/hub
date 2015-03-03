package commands

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

type listFlag []string

func (l *listFlag) String() string {
	return strings.Join([]string(*l), ",")
}

func (l *listFlag) Set(value string) error {
	for _, flag := range strings.Split(value, ",") {
		*l = append(*l, flag)
	}
	return nil
}

func isDir(file string) bool {
	f, err := os.Open(file)
	if err != nil {
		return false
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return false
	}

	return fi.IsDir()
}

func gitRemoteForProject(project *github.Project) (foundRemote *github.Remote) {
	remotes, err := github.Remotes()
	utils.Check(err)
	for _, remote := range remotes {
		remoteProject, pErr := remote.Project()
		if pErr == nil && remoteProject.SameAs(project) {
			foundRemote = &remote
			return
		}
	}

	return nil
}

func isEmptyDir(path string) bool {
	fullPath := filepath.Join(path, "*")
	match, _ := filepath.Glob(fullPath)
	return match == nil
}

func getTitleAndBodyFromFlags(messageFlag, fileFlag string) (title, body string, err error) {
	if messageFlag != "" {
		title, body = readMsg(messageFlag)
	} else if fileFlag != "" {
		var (
			content []byte
			err     error
		)

		if fileFlag == "-" {
			content, err = ioutil.ReadAll(os.Stdin)
		} else {
			content, err = ioutil.ReadFile(fileFlag)
		}
		utils.Check(err)

		title, body = readMsg(string(content))
	}

	return
}

func readMsg(msg string) (title, body string) {
	s := bufio.NewScanner(strings.NewReader(msg))
	if s.Scan() {
		title = s.Text()
		body = strings.TrimLeft(msg, title)

		title = strings.TrimSpace(title)
		body = strings.TrimSpace(body)
	}

	return
}

func runInLocalRepo(fn func(localRepo *github.GitHubRepo, project *github.Project, client *github.Client)) {
	localRepo, err := github.LocalRepo()
	utils.Check(err)

	project, err := localRepo.CurrentProject()
	utils.Check(err)

	client := github.NewClient(project.Host)
	fn(localRepo, project, client)

	os.Exit(0)
}
