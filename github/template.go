package github

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	PullRequestTemplate = "pull_request_template"
	IssueTemplate       = "issue_template"
	githubTemplateDir   = ".github"
)

func ReadTemplate(kind, workdir string) (body string, err error) {
	templateDir := filepath.Join(workdir, githubTemplateDir)

	path, err := getFilePath(templateDir, kind)
	if err != nil || path == "" {
		path, err = getFilePath(workdir, kind)
	}

	if path != "" {
		body, err = readContentsFromFile(path)
	}
	return
}

func getFilePath(dir, pattern string) (found string, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		fileName := file.Name()
		path := fileName

		if ext := filepath.Ext(fileName); ext == ".md" {
			path = strings.TrimRight(fileName, ".md")
		} else if ext == ".txt" {
			path = strings.TrimRight(fileName, ".txt")
		}

		path = strings.ToLower(path)

		if ok, _ := filepath.Match(pattern, path); ok {
			found = filepath.Join(dir, fileName)
			return
		}
	}
	return
}

func readContentsFromFile(filename string) (contents string, err error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	contents = strings.Replace(string(content), "\r\n", "\n", -1)
	contents = strings.TrimSpace(contents)
	return
}
