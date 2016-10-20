package github

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/github/hub/git"
)

func LocalRepo() (repo *GitHubRepo, err error) {
	repo = &GitHubRepo{}

	_, err = git.Dir()
	if err != nil {
		err = fmt.Errorf("fatal: Not a git repository")
		return
	}

	return
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
	remotesMap := make(map[string]Remote)

	if owner != "" {
		for _, remote := range r.remotes {
			p, e := remote.Project()
			if e == nil && strings.EqualFold(p.Owner, owner) {
				remotesMap[remote.Name] = remote
			}
		}
	}

	names := OriginNamesInLookupOrder
	for _, name := range names {
		if _, ok := remotesMap[name]; ok {
			continue
		}

		remote, err := r.RemoteByName(name)
		if err == nil {
			remotesMap[remote.Name] = *remote
		}
	}

	for i := len(names) - 1; i >= 0; i-- {
		name := names[i]
		if remote, ok := remotesMap[name]; ok {
			remotes = append(remotes, remote)
			delete(remotesMap, name)
		}
	}

	// anything other than names has higher priority
	for _, remote := range remotesMap {
		remotes = append([]Remote{remote}, remotes...)
	}

	return
}

func (r *GitHubRepo) CurrentBranch() (branch *Branch, err error) {
	head, err := git.Head()
	if err != nil {
		err = fmt.Errorf("Aborted: not currently on any branch.")
		return
	}

	branch = &Branch{r, head}
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

	branch = &Branch{r, name}

	return
}

func (r *GitHubRepo) RemoteBranchAndProject(owner string, preferUpstream bool) (branch *Branch, project *Project, err error) {
	r.loadRemotes()

	for _, remote := range r.remotes {
		if p, err := remote.Project(); err == nil {
			project = p
			break
		}
	}

	branch, err = r.CurrentBranch()
	if err != nil {
		return
	}

	if project != nil {
		branch = branch.PushTarget(owner, preferUpstream)
		if branch != nil && branch.IsRemote() {
			remote, e := r.RemoteByName(branch.RemoteName())
			if e == nil {
				project, err = remote.Project()
				if err != nil {
					return
				}
			}
		}
	}

	return
}

func (r *GitHubRepo) RemoteForRepo(repo *Repository) (*Remote, error) {
	r.loadRemotes()

	repoUrl, err := url.Parse(repo.HtmlUrl)
	if err != nil {
		return nil, err
	}

	project := NewProject(repo.Owner.Login, repo.Name, repoUrl.Host)

	for _, remote := range r.remotes {
		if rp, err := remote.Project(); err == nil {
			if rp.SameAs(project) {
				return &remote, nil
			}
		}
	}

	return nil, fmt.Errorf("could not find git remote for %s/%s", repo.Owner.Login, repo.Name)
}

func (r *GitHubRepo) OriginRemote() (remote *Remote, err error) {
	return r.RemoteByName("origin")
}

func (r *GitHubRepo) MainRemote() (remote *Remote, err error) {
	r.loadRemotes()

	if len(r.remotes) > 0 {
		remote = &r.remotes[0]
	}

	if remote == nil {
		err = fmt.Errorf("Can't find git remote origin")
	}

	return
}

func (r *GitHubRepo) MainProject() (project *Project, err error) {
	origin, err := r.MainRemote()
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
