package github

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestParseNameAndOwner(t *testing.T) {
	owner, name := parseOwnerAndName()
	assert.Equal(t, "gh", name)
	assert.Equal(t, "jingweno", owner)
}

func TestMustMatchGitHubUrl(t *testing.T) {
	url, _ := mustMatchGitHubUrl("git://github.com/jingweno/gh.git")
	assert.Equal(t, "git://github.com/jingweno/gh.git", url[0])

	url, _ = mustMatchGitHubUrl("git@github.com:jingweno/gh.git")
	assert.Equal(t, "git@github.com:jingweno/gh.git", url[0])

	url, _ = mustMatchGitHubUrl("https://github.com/jingweno/gh.git")
	assert.Equal(t, "https://github.com/jingweno/gh.git", url[0])
}
