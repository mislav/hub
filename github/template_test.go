package github

import (
	"path/filepath"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/github/hub/fixtures"
)

var prContent = `Description
-----------
[Enter your pull request description here]
`

var issueContent = `Description
-----------
[Enter your issue description here]
`

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
			"prTempalte":    pullRequestTemplate + ".md",
			"issueTempalte": issueTemplate + ".md",
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
			"prTempalte":    pullRequestTemplate + ".md",
			"issueTempalte": issueTemplate + ".md",
			"dir":           githubTemplateDir,
		})

	assert.Equal(t, prContent, GetPullRequestTemplate())
	assert.Equal(t, issueContent, GetIssueTemplate())
}

// When no default message is provided, two blank lines should be added
// (representing the pull request title), and the left should be template.
func TestGeneratePRTemplate_NoDefaultMessage(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	addGithubTemplates(repo, map[string]string{})

	defaultMessage := ""
	expectedOutput := `

Description
-----------
[Enter your pull request description here]
`

	assert.Equal(t, expectedOutput, GeneratePRTemplate(defaultMessage))
}

// When a single line commit message is provided, the commit message should
// encompass the first line, then a empty new line, then the template.
func TestGeneratePRTemplate_SingleLineDefaultMessage(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	addGithubTemplates(repo, map[string]string{})

	defaultMessage := "Add Pull Request Templates to Hub"
	expectedOutput := `Add Pull Request Templates to Hub

Description
-----------
[Enter your pull request description here]
`

	assert.Equal(t, expectedOutput, GeneratePRTemplate(defaultMessage))
}

// When a multi line commit message is provided, the first line of the commit
// message should be the first line, then a empty new line, then the template.
//
// TODO (maybe):  Allow for templates to support auto filling the description
// section with the rest of the commit message.
func TestGeneratePRTemplate_MultiLineDefaultMessage(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	addGithubTemplates(repo, map[string]string{})

	defaultMessage := `Add Pull Request Templates to Hub

Allow repo maintainers to set a default template and allow developers to
continue to use hub!
`
	expectedOutput := `Add Pull Request Templates to Hub

Description
-----------
[Enter your pull request description here]
`

	assert.Equal(t, expectedOutput, GeneratePRTemplate(defaultMessage))
}

func addGithubTemplates(r *fixtures.TestRepo, config map[string]string) {
	repoDir := "test.git"
	if dir := config["dir"]; dir != "" {
		repoDir = filepath.Join(repoDir, dir)
	}

	prTemplatePath := filepath.Join(repoDir, pullRequestTemplate)
	if prTmplPath := config["prTemplate"]; prTmplPath != "" {
		prTemplatePath = prTmplPath
	}

	issueTemplatePath := filepath.Join(repoDir, issueTemplate)
	if issueTmplPath := config["issueTemplate"]; issueTmplPath != "" {
		issueTemplatePath = issueTmplPath
	}

	r.AddFile(prTemplatePath, prContent)
	r.AddFile(issueTemplatePath, issueContent)
}
