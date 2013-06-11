package github

import (
	"strings"
)

type Repo struct {
	Base    string
	Head    string
	Project *Project
}

func (r *Repo) FullBase() string {
	if strings.Contains(r.Base, ":") {
		return r.Base
	}

	return r.Project.Owner + ":" + r.Base
}

func (r *Repo) FullHead() string {
	if strings.Contains(r.Head, ":") {
		return r.Head
	}

	return r.Project.Owner + ":" + r.Head
}
