package github

import (
	"fmt"
	"github.com/jingweno/gh/git"
	"strings"
)

func LocalRepo() *GitHubRepo {
	return &GitHubRepo{}
}

type GitHubRepo struct {
	remotes []git.Remote
}

func (r *GitHubRepo) remotesByName(name string) (*git.Remote, error) {
	if r.remotes == nil {
		remotes, err := git.Remotes()
		if err != nil {
			return nil, err
		}
		r.remotes = remotes
	}

	for _, remote := range r.remotes {
		if remote.Name == name {
			return &remote, nil
		}
	}

	return nil, fmt.Errorf("No git remote with name %s", name)
}

func (r *GitHubRepo) CurrentBranch() (branch *Branch, err error) {
	head, err := git.Head()
	if err != nil {
		err = fmt.Errorf("Aborted: not currently on any branch.")
		return
	}

	branch = &Branch{head}
	return
}

func (r *GitHubRepo) MasterBranch() (branch *Branch, err error) {
	origin, err := r.remotesByName("origin")
	if err != nil {
		return
	}

	name, err := git.SymbolicFullName(origin.Name)
	if err != nil {
		name = "refs/head/master"
		err = nil
	}

	branch = &Branch{name}

	return
}

func (r *GitHubRepo) MainProject() (project *Project, err error) {
	origin, err := r.remotesByName("origin")
	if err != nil {
		err = fmt.Errorf("Aborted: the origin remote doesn't point to a GitHub repository.")
		return
	}

	project, err = NewProjectFromURL(origin.URL)
	if err != nil {
		err = fmt.Errorf("Aborted: the origin remote doesn't point to a GitHub repository.")
	}

	return
}

func (r *GitHubRepo) CurrentProject() (project *Project, err error) {
	project, err = r.UpstreamProject()
	if err != nil {
		project, err = r.MainProject()
	}

	return
}

func (r *GitHubRepo) UpstreamProject() (project *Project, err error) {
	currentBranch, err := r.CurrentBranch()
	if err != nil {
		return
	}

	upstream, err := currentBranch.Upstream()
	if err != nil {
		return
	}

	remote, err := r.remotesByName(upstream.RemoteName())
	if err != nil {
		return
	}

	project, err = NewProjectFromURL(remote.URL)

	return
}

// TODO: remove it
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
