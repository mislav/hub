package main

import (
	"strings"
)

type Repo struct {
	Base    string
	Head    string
	Project *GitHubProject
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

	project := CurrentProject()

	return &Repo{base, head, project}
}
