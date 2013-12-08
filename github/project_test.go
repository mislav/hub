package github

import (
	"github.com/bmizerany/assert"
	"net/url"
	"testing"
)

func TestNewProjectOwnerAndName(t *testing.T) {
	CreateTestConfig("jingweno", "123")

	project := NewProjectFromOwnerAndName("", "mojombo/gh")
	assert.Equal(t, "mojombo", project.Owner)
	assert.Equal(t, "gh", project.Name)

	project = NewProjectFromOwnerAndName("progmob", "mojombo/gh")
	assert.Equal(t, "progmob", project.Owner)
	assert.Equal(t, "gh", project.Name)

	project = NewProjectFromOwnerAndName("mojombo/jekyll", "gh")
	assert.Equal(t, "mojombo", project.Owner)
	assert.Equal(t, "gh", project.Name)

	project = NewProjectFromOwnerAndName("mojombo/gh", "")
	assert.Equal(t, "mojombo", project.Owner)
	assert.Equal(t, "gh", project.Name)

	project = NewProjectFromOwnerAndName("", "gh")
	assert.Equal(t, "jingweno", project.Owner)
	assert.Equal(t, "gh", project.Name)

	project = NewProjectFromOwnerAndName("", "jingweno/gh/foo")
	assert.Equal(t, "jingweno", project.Owner)
	assert.Equal(t, "gh/foo", project.Name)

	project = NewProjectFromOwnerAndName("mojombo", "gh")
	assert.Equal(t, "mojombo", project.Owner)
	assert.Equal(t, "gh", project.Name)
}

func TestWebURL(t *testing.T) {
	project := Project{Name: "foo", Owner: "bar"}
	url := project.WebURL("", "", "baz")
	assert.Equal(t, "https://github.com/bar/foo/baz", url)

	url = project.WebURL("1", "2", "")
	assert.Equal(t, "https://github.com/2/1", url)
}

func TestGitURL(t *testing.T) {
	project := Project{Name: "foo", Owner: "bar"}
	url := project.GitURL("gh", "jingweno", false)
	assert.Equal(t, "git://github.com/jingweno/gh.git", url)

	url = project.GitURL("gh", "jingweno", true)
	assert.Equal(t, "git@github.com:jingweno/gh.git", url)
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
	u, _ := url.Parse("ssh://git@github.com/jingweno/gh.git")
	p, err := NewProjectFromURL(u)

	assert.Equal(t, nil, err)
	assert.Equal(t, "gh", p.Name)
	assert.Equal(t, "jingweno", p.Owner)

	u, _ = url.Parse("git://github.com/jingweno/gh.git")
	p, err = NewProjectFromURL(u)

	assert.Equal(t, nil, err)
	assert.Equal(t, "gh", p.Name)
	assert.Equal(t, "jingweno", p.Owner)

	u, _ = url.Parse("https://github.com/jingweno/gh")
	p, err = NewProjectFromURL(u)

	assert.Equal(t, nil, err)
	assert.Equal(t, "gh", p.Name)
	assert.Equal(t, "jingweno", p.Owner)
}
