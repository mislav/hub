package main

import (
	"strings"
)

type Repo struct {
	Dir     string
	Editor  string
	Owner   string
	Project string
	Base    string
	Head    string
}

func (r *Repo) FullBase() string {
	if strings.Contains(r.Base, ":") {
		return r.Base
	} else {
		return r.Owner + ":" + r.Base
	}
}

func (r *Repo) FullHead() string {
	if strings.Contains(r.Head, ":") {
		return r.Head
	} else {
		return r.Owner + ":" + r.Head
	}
}

func NewRepo(base, head string) *Repo {
	if base == "" {
		base = "master"
	}
	if head == "" {
		head, _ = FetchGitHead()
	}

	dir, _ := FetchGitDir()
	editor, _ := FetchGitEditor()
	owner, _ := FetchGitOwner()
	project, _ := FetchGitProject()

	return &Repo{dir, editor, owner, project, base, head}
}
