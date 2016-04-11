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
