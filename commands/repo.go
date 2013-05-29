package commands

import (
	"github.com/jingweno/gh/git"
	"github.com/jingweno/gh/github"
	"strings"
)

type Repo struct {
	Base    string
	Head    string
	Project *github.Project
}

func (r *Repo) FullBase() string {
	if strings.Contains(r.Base, ":") {
		return r.Base
	} else {
		return r.Project.Owner + ":" + r.Base
	}
}

func (r *Repo) FullHead() string {
	if strings.Contains(r.Head, ":") {
		return r.Head
	} else {
		return r.Project.Owner + ":" + r.Head
	}
}

func NewRepo(base, head string) *Repo {
	if base == "" {
		base = "master"
	}
	if head == "" {
		head, _ = git.Head()
	}

	project := github.CurrentProject()

	return &Repo{base, head, project}
}
