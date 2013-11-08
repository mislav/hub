package commands

import (
	"github.com/bmizerany/assert"
	"github.com/octokit/go-octokit/octokit"
	"testing"
)

func TestFetchAndMerge(t *testing.T) {
	url := "https://github.com/jingweno/gh/pull/73"
	number := 73
	title := "title"

	args := NewArgs([]string{"merge", url})

	userLogin := "jingweno"
	user := octokit.User{Login: userLogin}

	repoPrivate := true
	repo := octokit.Repository{Private: repoPrivate}

	headRef := "new-feature"
	head := octokit.Commit{Ref: headRef, Repo: repo, Label: "jingweno:new-feature"}

	pullRequest := octokit.PullRequest{Number: number, Title: title, HTMLURL: url, User: user, Head: head}

	err := fetchAndMerge(args, &pullRequest)
	assert.Equal(t, nil, err)

	cmds := args.Commands()
	assert.Equal(t, 2, len(cmds))

	cmd := cmds[0]
	assert.Equal(t, "git fetch git@github.com:jingweno/gh.git +refs/heads/new-feature:refs/remotes/jingweno/new-feature", cmd.String())

	cmd = cmds[1]
	assert.Equal(t, "git merge jingweno/new-feature --no-ff -m Merge pull request #73 from jingweno/new-feature\n\ntitle", cmd.String())
}
