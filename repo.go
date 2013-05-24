package main

type Repo struct {
	Dir     string
	Editor  string
	Owner   string
	Project string
	Base    string
	Head    string
}

func (r *Repo) FullBase() string {
	return r.Owner + ":" + r.Base
}

func (r *Repo) FullHead() string {
	return r.Owner + ":" + r.Head
}

func NewRepo() *Repo {
	dir, _ := FetchGitDir()
	editor, _ := FetchGitEditor()
	owner, _ := FetchGitOwner()
	project, _ := FetchGitProject()
	head, _ := FetchGitHead()

	return &Repo{dir, editor, owner, project, "master", head}
}
