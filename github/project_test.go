package github

import (
	"net/url"
	"os"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/github/hub/fixtures"
)

func TestProject_WebURL(t *testing.T) {
	project := Project{
		Name:     "foo",
		Owner:    "bar",
		Host:     "github.com",
		Protocol: "https",
	}

	myUrl := project.WebURL("", "", "baz")
	assert.Equal(t, "https://github.com/bar/foo/baz", myUrl)

	myUrl = project.WebURL("1", "2", "")
	assert.Equal(t, "https://github.com/2/1", myUrl)

	myUrl = project.WebURL("hub.wiki", "defunkt", "")
	assert.Equal(t, "https://github.com/defunkt/hub/wiki", myUrl)

	myUrl = project.WebURL("hub.wiki", "defunkt", "commits")
	assert.Equal(t, "https://github.com/defunkt/hub/wiki/_history", myUrl)

	myUrl = project.WebURL("hub.wiki", "defunkt", "pages")
	assert.Equal(t, "https://github.com/defunkt/hub/wiki/_pages", myUrl)
}

func TestProject_GitURL(t *testing.T) {
	os.Setenv("HUB_PROTOCOL", "https")
	defer os.Setenv("HUB_PROTOCOL", "")

	project := Project{
		Name:  "foo",
		Owner: "bar",
		Host:  "github.com",
	}

	myUrl := project.GitURL("gh", "jingweno", false)
	assert.Equal(t, "https://github.com/jingweno/gh.git", myUrl)

	os.Setenv("HUB_PROTOCOL", "git")
	myUrl = project.GitURL("gh", "jingweno", false)
	assert.Equal(t, "git://github.com/jingweno/gh.git", myUrl)

	os.Setenv("HUB_PROTOCOL", "ssh")
	myUrl = project.GitURL("gh", "jingweno", true)
	assert.Equal(t, "git@github.com:jingweno/gh.git", myUrl)

	myUrl = project.GitURL("gh", "jingweno", true)
	assert.Equal(t, "git@github.com:jingweno/gh.git", myUrl)
}

func TestProject_GitURLEnterprise(t *testing.T) {
	project := Project{
		Name:  "foo",
		Owner: "bar",
		Host:  "https://github.corporate.com",
	}

	defer os.Setenv("HUB_PROTOCOL", "")

	os.Setenv("HUB_PROTOCOL", "https")
	myUrl := project.GitURL("gh", "jingweno", false)
	assert.Equal(t, "https://github.corporate.com/jingweno/gh.git", myUrl)

	os.Setenv("HUB_PROTOCOL", "ssh")
	myUrl = project.GitURL("gh", "jingweno", false)
	assert.Equal(t, "git@github.corporate.com:jingweno/gh.git", myUrl)

	os.Setenv("HUB_PROTOCOL", "git")
	myUrl = project.GitURL("gh", "jingweno", false)
	assert.Equal(t, "git://github.corporate.com/jingweno/gh.git", myUrl)

	myUrl = project.GitURL("gh", "jingweno", true)
	assert.Equal(t, "git@github.corporate.com:jingweno/gh.git", myUrl)
}

func TestProject_NewProjectFromURL(t *testing.T) {
	testConfigs := fixtures.SetupTestConfigs()
	defer testConfigs.TearDown()

	u, _ := url.Parse("ssh://git@github.com/octokit/go-octokit.git")
	p, err := NewProjectFromURL(u)

	assert.Equal(t, nil, err)
	assert.Equal(t, "go-octokit", p.Name)
	assert.Equal(t, "octokit", p.Owner)
	assert.Equal(t, "github.com", p.Host)
	assert.Equal(t, "http", p.Protocol)

	u, _ = url.Parse("ssh://ssh.github.com/octokit/go-octokit.git")
	p, err = NewProjectFromURL(u)

	assert.Equal(t, nil, err)
	assert.Equal(t, "go-octokit", p.Name)
	assert.Equal(t, "octokit", p.Owner)
	assert.Equal(t, "github.com", p.Host)
	assert.Equal(t, "http", p.Protocol)

	u, _ = url.Parse("git://github.com/octokit/go-octokit.git")
	p, err = NewProjectFromURL(u)

	assert.Equal(t, nil, err)
	assert.Equal(t, "go-octokit", p.Name)
	assert.Equal(t, "octokit", p.Owner)
	assert.Equal(t, "github.com", p.Host)
	assert.Equal(t, "http", p.Protocol)

	u, _ = url.Parse("https://github.com/octokit/go-octokit")
	p, err = NewProjectFromURL(u)

	assert.Equal(t, nil, err)
	assert.Equal(t, "go-octokit", p.Name)
	assert.Equal(t, "octokit", p.Owner)
	assert.Equal(t, "github.com", p.Host)
	assert.Equal(t, "https", p.Protocol)

	u, _ = url.Parse("origin/master")
	_, err = NewProjectFromURL(u)

	assert.NotEqual(t, nil, err)
}
