package github

import (
	"github.com/bmizerany/assert"
	"net/url"
	"os"
	"testing"
)

func TestWebURL(t *testing.T) {
	project := Project{Name: "foo", Owner: "bar", Host: "github.com"}
	url := project.WebURL("", "", "baz")
	assert.Equal(t, "https://github.com/bar/foo/baz", url)

	url = project.WebURL("1", "2", "")
	assert.Equal(t, "https://github.com/2/1", url)

	url = project.WebURL("hub.wiki", "defunkt", "")
	assert.Equal(t, "https://github.com/defunkt/hub/wiki", url)

	url = project.WebURL("hub.wiki", "defunkt", "commits")
	assert.Equal(t, "https://github.com/defunkt/hub/wiki/_history", url)

	url = project.WebURL("hub.wiki", "defunkt", "pages")
	assert.Equal(t, "https://github.com/defunkt/hub/wiki/_pages", url)
}

func TestGitURLGitHub(t *testing.T) {
	os.Setenv("GH_PROTOCOL", "https")
	project := Project{Name: "foo", Owner: "bar", Host: "github.com"}

	url := project.GitURL("gh", "jingweno", false)
	assert.Equal(t, "https://github.com/jingweno/gh.git", url)

	os.Setenv("GH_PROTOCOL", "git")
	url = project.GitURL("gh", "jingweno", false)
	assert.Equal(t, "git://github.com/jingweno/gh.git", url)

	url = project.GitURL("gh", "jingweno", true)
	assert.Equal(t, "git@github.com:jingweno/gh.git", url)
}

func TestGitURLEnterprise(t *testing.T) {
	project := Project{Name: "foo", Owner: "bar", Host: "https://github.corporate.com"}

	os.Setenv("GH_PROTOCOL", "https")
	url := project.GitURL("gh", "jingweno", false)
	assert.Equal(t, "https://github.corporate.com/jingweno/gh.git", url)

	os.Setenv("GH_PROTOCOL", "git")
	url = project.GitURL("gh", "jingweno", false)
	assert.Equal(t, "git://github.corporate.com/jingweno/gh.git", url)

	url = project.GitURL("gh", "jingweno", true)
	assert.Equal(t, "git@github.corporate.com:jingweno/gh.git", url)
}

func TestParseOwnerAndName(t *testing.T) {
	owner, name := parseOwnerAndName("git://github.com/jingweno/gh.git")
	assert.Equal(t, "gh", name)
	assert.Equal(t, "jingweno", owner)
}

func TestMustMatchGitHubURL(t *testing.T) {
	url, _ := mustMatchGitHubURL("git://github.com/jingweno/gh.git")
	assert.Equal(t, "git://github.com/jingweno/gh.git", url[0])
	assert.Equal(t, "jingweno", url[1])
	assert.Equal(t, "gh", url[2])

	url, _ = mustMatchGitHubURL("git://github.com/jingweno/gh")
	assert.Equal(t, "git://github.com/jingweno/gh", url[0])
	assert.Equal(t, "jingweno", url[1])
	assert.Equal(t, "gh", url[2])

	url, _ = mustMatchGitHubURL("git://git@github.com/jingweno/gh")
	assert.Equal(t, "git://git@github.com/jingweno/gh", url[0])
	assert.Equal(t, "jingweno", url[1])
	assert.Equal(t, "gh", url[2])

	url, _ = mustMatchGitHubURL("git@github.com:jingweno/gh.git")
	assert.Equal(t, "git@github.com:jingweno/gh.git", url[0])
	assert.Equal(t, "jingweno", url[1])
	assert.Equal(t, "gh", url[2])

	url, _ = mustMatchGitHubURL("git@github.com:jingweno/gh")
	assert.Equal(t, "git@github.com:jingweno/gh", url[0])
	assert.Equal(t, "jingweno", url[1])
	assert.Equal(t, "gh", url[2])

	url, _ = mustMatchGitHubURL("https://github.com/jingweno/gh.git")
	assert.Equal(t, "https://github.com/jingweno/gh.git", url[0])
	assert.Equal(t, "jingweno", url[1])
	assert.Equal(t, "gh", url[2])

	url, _ = mustMatchGitHubURL("https://github.com/jingweno/gh")
	assert.Equal(t, "https://github.com/jingweno/gh", url[0])
	assert.Equal(t, "jingweno", url[1])
	assert.Equal(t, "gh", url[2])
}

func TestNewProjectFromURL(t *testing.T) {
	u, _ := url.Parse("ssh://git@github.com/octokit/go-octokit.git")
	p, err := NewProjectFromURL(u)

	assert.Equal(t, nil, err)
	assert.Equal(t, "go-octokit", p.Name)
	assert.Equal(t, "octokit", p.Owner)

	u, _ = url.Parse("git://github.com/octokit/go-octokit.git")
	p, err = NewProjectFromURL(u)

	assert.Equal(t, nil, err)
	assert.Equal(t, "go-octokit", p.Name)
	assert.Equal(t, "octokit", p.Owner)

	u, _ = url.Parse("https://github.com/octokit/go-octokit")
	p, err = NewProjectFromURL(u)

	assert.Equal(t, nil, err)
	assert.Equal(t, "go-octokit", p.Name)
	assert.Equal(t, "octokit", p.Owner)

	u, _ = url.Parse("origin/master")
	_, err = NewProjectFromURL(u)

	assert.NotEqual(t, nil, err)
}
