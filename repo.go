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
	dir := FetchGitDir()
	editor := FetchGitEditor()
	owner := FetchGitOwner()
	project := FetchGitProject()
	head := FetchGitHead()

	return &Repo{dir, editor, owner, project, "master", head}
}
