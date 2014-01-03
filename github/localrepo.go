package github

import (
	"fmt"
	"github.com/jingweno/gh/git"
)

func LocalRepo() *GitHubRepo {
	return &GitHubRepo{}
}

type GitHubRepo struct {
	remotes []Remote
}

func (r *GitHubRepo) loadRemotes() error {
	if r.remotes != nil {
		return nil
	}

	remotes, err := Remotes()
	if err != nil {
		return err
	}
	r.remotes = remotes

	return nil
}

func (r *GitHubRepo) RemoteByName(name string) (*Remote, error) {
	r.loadRemotes()

	for _, remote := range r.remotes {
		if remote.Name == name {
			return &remote, nil
		}
	}

	return nil, fmt.Errorf("No git remote with name %s", name)
}

func (r *GitHubRepo) remotesForPublish(owner string) (remotes []Remote) {
	r.loadRemotes()

	if owner != "" {
		for _, remote := range r.remotes {
			p, e := remote.Project()
			if e == nil && p.Owner == owner {
				remotes = append(remotes, remote)
			}
		}
	}

	remote, err := r.RemoteByName("origin")
	if err == nil {
		remotes = append(remotes, *remote)
	}

	remote, err = r.RemoteByName("github")
	if err == nil {
		remotes = append(remotes, *remote)
	}

	remote, err = r.RemoteByName("upstream")
	if err == nil {
		remotes = append(remotes, *remote)
	}

	return
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

func (r *GitHubRepo) MasterBranch() (branch *Branch) {
	origin, e := r.RemoteByName("origin")
	var name string
	if e == nil {
		name, _ = git.BranchAtRef("refs", "remotes", origin.Name, "HEAD")
	}

	if name == "" {
		name = "refs/heads/master"
	}

	branch = &Branch{name}

	return
}

func (r *GitHubRepo) RemoteBranchAndProject(owner string) (branch *Branch, project *Project, err error) {
	project, err = r.MainProject()
	if err != nil {
		return
	}

	branch, err = r.CurrentBranch()
	if err != nil {
		return
	}

	pushDefault, _ := git.Config("push.default")
	if pushDefault == "upstream" || pushDefault == "tracking" {
		branch, err = branch.Upstream()
		if err != nil {
			return
		}
	} else {
		shortName := branch.ShortName()
		remotes := r.remotesForPublish(owner)
		for _, remote := range remotes {
			if git.HasFile("refs", "remotes", remote.Name, shortName) {
				name := fmt.Sprintf("refs/remotes/%s/%s", remote.Name, shortName)
				branch = &Branch{name}
				break
			}
		}
	}

	if branch.IsRemote() {
		remote, e := r.RemoteByName(branch.RemoteName())
		if e == nil {
			project, err = remote.Project()
			if err != nil {
				return
			}
		}
	}

	return
}

func (r *GitHubRepo) MainProject() (project *Project, err error) {
	origin, err := r.RemoteByName("origin")
	if err != nil {
		err = fmt.Errorf("Aborted: the origin remote doesn't point to a GitHub repository.")
		return
	}

	project, err = origin.Project()
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

	remote, err := r.RemoteByName(upstream.RemoteName())
	if err != nil {
		return
	}

	project, err = remote.Project()

	return
}
