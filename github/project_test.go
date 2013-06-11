package github

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestWebURL(t *testing.T) {
	project := Project{"foo", "bar"}
	url := project.WebURL("", "", "baz")
	assert.Equal(t, "https://github.com/bar/foo/baz", url)

	url = project.WebURL("1", "2", "")
	assert.Equal(t, "https://github.com/2/1", url)
}

func TestParseNameAndOwner(t *testing.T) {
	owner, name := parseOwnerAndName()
	assert.Equal(t, "gh", name)
	assert.Equal(t, "jingweno", owner)
}

func TestMustMatchGitHubURL(t *testing.T) {
	url, _ := mustMatchGitHubURL("git://github.com/jingweno/gh.git")
	assert.Equal(t, "git://github.com/jingweno/gh.git", url[0])

	url, _ = mustMatchGitHubURL("git@github.com:jingweno/gh.git")
	assert.Equal(t, "git@github.com:jingweno/gh.git", url[0])

	url, _ = mustMatchGitHubURL("https://github.com/jingweno/gh.git")
	assert.Equal(t, "https://github.com/jingweno/gh.git", url[0])
}
