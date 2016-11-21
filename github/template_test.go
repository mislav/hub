package github

import (
	"os"
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

	pwd, _ := os.Getwd()
	tpl, err := ReadTemplate(PullRequestTemplate, pwd)
	assert.Equal(t, nil, err)
	assert.Equal(t, "", tpl)

	tpl, err = ReadTemplate(IssueTemplate, pwd)
	assert.Equal(t, nil, err)
	assert.Equal(t, "", tpl)
}

func TestGithubTemplate_withInvalidTemplate(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	addGithubTemplates(repo, map[string]string{"dir": "invalidPath"})

	pwd, _ := os.Getwd()
	tpl, err := ReadTemplate(PullRequestTemplate, pwd)
	assert.Equal(t, nil, err)
	assert.Equal(t, "", tpl)

	tpl, err = ReadTemplate(IssueTemplate, pwd)
	assert.Equal(t, nil, err)
	assert.Equal(t, "", tpl)
}

func TestGithubTemplate_WithMarkdown(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	addGithubTemplates(repo,
		map[string]string{
			"prTemplate":    PullRequestTemplate + ".md",
			"issueTemplate": IssueTemplate + ".md",
		})

	pwd, _ := os.Getwd()
	tpl, err := ReadTemplate(PullRequestTemplate, pwd)
	assert.Equal(t, nil, err)
	assert.Equal(t, prContent, tpl)

	tpl, err = ReadTemplate(IssueTemplate, pwd)
	assert.Equal(t, nil, err)
	assert.Equal(t, issueContent, tpl)
}

func TestGithubTemplate_WithTemplateInHome(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	addGithubTemplates(repo, map[string]string{})

	pwd, _ := os.Getwd()
	tpl, err := ReadTemplate(PullRequestTemplate, pwd)
	assert.Equal(t, nil, err)
	assert.Equal(t, prContent, tpl)

	tpl, err = ReadTemplate(IssueTemplate, pwd)
	assert.Equal(t, nil, err)
	assert.Equal(t, issueContent, tpl)
}

func TestGithubTemplate_WithTemplateInGithubDir(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	addGithubTemplates(repo, map[string]string{"dir": githubTemplateDir})

	pwd, _ := os.Getwd()
	tpl, err := ReadTemplate(PullRequestTemplate, pwd)
	assert.Equal(t, nil, err)
	assert.Equal(t, prContent, tpl)

	tpl, err = ReadTemplate(IssueTemplate, pwd)
	assert.Equal(t, nil, err)
	assert.Equal(t, issueContent, tpl)
}

func TestGithubTemplate_WithTemplateInGithubDirAndMarkdown(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	addGithubTemplates(repo,
		map[string]string{
			"prTemplate":    PullRequestTemplate + ".md",
			"issueTemplate": IssueTemplate + ".md",
			"dir":           githubTemplateDir,
		})

	pwd, _ := os.Getwd()
	tpl, err := ReadTemplate(PullRequestTemplate, pwd)
	assert.Equal(t, nil, err)
	assert.Equal(t, prContent, tpl)

	tpl, err = ReadTemplate(IssueTemplate, pwd)
	assert.Equal(t, nil, err)
	assert.Equal(t, issueContent, tpl)
}

func addGithubTemplates(r *fixtures.TestRepo, config map[string]string) {
	repoDir := "test.git"
	if dir := config["dir"]; dir != "" {
		repoDir = filepath.Join(repoDir, dir)
	}

	prTemplatePath := filepath.Join(repoDir, PullRequestTemplate)
	if prTmplPath := config["prTemplate"]; prTmplPath != "" {
		prTemplatePath = filepath.Join(repoDir, prTmplPath)
	}

	issueTemplatePath := filepath.Join(repoDir, IssueTemplate)
	if issueTmplPath := config["issueTemplate"]; issueTmplPath != "" {
		issueTemplatePath = filepath.Join(repoDir, issueTmplPath)
	}

	r.AddFile(prTemplatePath, prContent)
	r.AddFile(issueTemplatePath, issueContent)
}
