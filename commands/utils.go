package commands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/github/hub/git"
	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
)

type stringSliceValue []string

func (s *stringSliceValue) Set(val string) error {
	*s = append(*s, val)
	return nil
}

func (s *stringSliceValue) String() string {
	return fmt.Sprintf("%s", *s)
}

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

func isCloneable(file string) bool {
	f, err := os.Open(file)
	if err != nil {
		return false
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return false
	}

	if fi.IsDir() {
		gitf, err := os.Open(filepath.Join(file, ".git"))
		if err == nil {
			gitf.Close()
			return true
		} else {
			return git.IsGitDir(file)
		}
	} else {
		reader := bufio.NewReader(f)
		line, err := reader.ReadString('\n')
		if err == nil {
			return strings.Contains(line, "git bundle")
		} else {
			return false
		}
	}
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

func readMsgFromFile(filename string, edit bool, editorPrefix, editorTopic string) (title, body string, editor *github.Editor, err error) {
	message, err := msgFromFile(filename)
	if err != nil {
		return
	}

	if edit {
		editor, err = github.NewEditor(editorPrefix, editorTopic, message)
		if err != nil {
			return
		}
		title, body, err = editor.EditTitleAndBody()
		return
	} else {
		title, body = readMsg(message)
		return
	}
}

func msgFromFile(filename string) (string, error) {
	var content []byte
	var err error

	if filename == "-" {
		content, err = ioutil.ReadAll(os.Stdin)
	} else {
		content, err = ioutil.ReadFile(filename)
	}
	if err != nil {
		return "", err
	}

	return strings.Replace(string(content), "\r\n", "\n", -1), nil
}

func readMsg(message string) (title, body string) {
	parts := strings.SplitN(message, "\n\n", 2)

	title = strings.TrimSpace(strings.Replace(parts[0], "\n", " ", -1))
	if len(parts) > 1 {
		body = strings.TrimSpace(parts[1])
	}
	return
}

func printBrowseOrCopy(args *Args, msg string, openBrowser bool, performCopy bool) {
	if performCopy {
		if err := clipboard.WriteAll(msg); err != nil {
			ui.Errorf("Error copying %s to clipboard:\n%s\n", msg, err.Error())
		}
	}

	if openBrowser {
		launcher, err := utils.BrowserLauncher()
		utils.Check(err)
		args.Replace(launcher[0], "", launcher[1:]...)
		args.AppendParams(msg)
	}

	if !openBrowser && !performCopy {
		args.AfterFn(func() error {
			ui.Println(msg)
			return nil
		})
	}
}
