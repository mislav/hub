package commands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/github/hub/v2/git"
	"github.com/github/hub/v2/ui"
	"github.com/github/hub/v2/utils"
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

type messageBlocks []string

func (m *messageBlocks) String() string {
	return strings.Join([]string(*m), "\n\n")
}

func (m *messageBlocks) Set(value string) error {
	*m = append(*m, value)
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
		}
		return git.IsGitDir(file)
	}
	reader := bufio.NewReader(f)
	line, err := reader.ReadString('\n')
	if err == nil {
		return strings.Contains(line, "git bundle")
	}
	return false
}

func isEmptyDir(path string) bool {
	fullPath := filepath.Join(path, "*")
	match, _ := filepath.Glob(fullPath)
	return match == nil
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
