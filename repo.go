package main

import (
	"strings"
)

type Repo struct {
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
		head, _ = git.Head()
	}

	owner, _ := git.Owner()
	project, _ := git.Project()

	return &Repo{owner, project, base, head}
}
