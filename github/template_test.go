package github

import (
	"path/filepath"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/github/hub/fixtures"
)

var prContent = `Description
-----------
[Enter your pull request description here]`

var issueContent = `Description
-----------
[Enter your issue description here]`

func TestGithubTemplate_withoutTemplate(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	assert.Equal(t, "", GetPullRequestTemplate())
	assert.Equal(t, "", GetIssueTemplate())
}

func TestGithubTemplate_withInvalidTemplate(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	addGithubTemplates(repo, map[string]string{"dir": "invalidPath"})

	assert.Equal(t, "", GetPullRequestTemplate())
	assert.Equal(t, "", GetIssueTemplate())
}

func TestGithubTemplate_WithMarkdown(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	addGithubTemplates(repo,
		map[string]string{
			"prTemplate":    pullRequestTemplate + ".md",
			"issueTemplate": issueTemplate + ".md",
		})

	assert.Equal(t, prContent, GetPullRequestTemplate())
	assert.Equal(t, issueContent, GetIssueTemplate())
}

func TestGithubTemplate_WithTemplateInHome(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	addGithubTemplates(repo, map[string]string{})

	assert.Equal(t, prContent, GetPullRequestTemplate())
	assert.Equal(t, issueContent, GetIssueTemplate())
}

func TestGithubTemplate_WithTemplateInGithubDir(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	addGithubTemplates(repo, map[string]string{"dir": githubTemplateDir})

	assert.Equal(t, prContent, GetPullRequestTemplate())
	assert.Equal(t, issueContent, GetIssueTemplate())
}

func TestGithubTemplate_WithTemplateInGithubDirAndMarkdown(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	addGithubTemplates(repo,
		map[string]string{
			"prTemplate":    pullRequestTemplate + ".md",
			"issueTemplate": issueTemplate + ".md",
			"dir":           githubTemplateDir,
		})

	assert.Equal(t, prContent, GetPullRequestTemplate())
	assert.Equal(t, issueContent, GetIssueTemplate())
}

func addGithubTemplates(r *fixtures.TestRepo, config map[string]string) {
	repoDir := "test.git"
	if dir := config["dir"]; dir != "" {
		repoDir = filepath.Join(repoDir, dir)
	}

	prTemplatePath := filepath.Join(repoDir, pullRequestTemplate)
	if prTmplPath := config["prTemplate"]; prTmplPath != "" {
		prTemplatePath = filepath.Join(repoDir, prTmplPath)
	}

	issueTemplatePath := filepath.Join(repoDir, issueTemplate)
	if issueTmplPath := config["issueTemplate"]; issueTmplPath != "" {
		issueTemplatePath = filepath.Join(repoDir, issueTmplPath)
	}

	r.AddFile(prTemplatePath, prContent)
	r.AddFile(issueTemplatePath, issueContent)
}
