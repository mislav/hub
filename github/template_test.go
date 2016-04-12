package github

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/github/hub/fixtures"
)

// When no default message is provided, two blank lines should be added
// (representing the pull request title), and the left should be template.
func TestGeneratePRTemplate_NoDefaultMessage(t *testing.T) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	repo.AddGithubTemplatesDir()

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

	repo.AddGithubTemplatesDir()

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

	repo.AddGithubTemplatesDir()

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
