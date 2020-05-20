package github

import (
	"net/url"
	"os"
	"testing"

	"github.com/github/hub/v2/fixtures"
	"github.com/github/hub/v2/internal/assert"
)

func TestSameAs(t *testing.T) {
	tests := []struct {
		name     string
		project1 *Project
		project2 *Project
		want     bool
	}{
		{
			name: "same project",
			project1: &Project{
				Owner: "fOo",
				Name:  "baR",
				Host:  "gItHUb.com",
			},
			project2: &Project{
				Owner: "FoO",
				Name:  "BAr",
				Host:  "GithUB.com",
			},
			want: true,
		},
		{
			name: "different project",
			project1: &Project{
				Owner: "foo",
				Name:  "bar",
				Host:  "github.com",
			},
			project2: &Project{
				Owner: "foo",
				Name:  "baz",
				Host:  "github.com",
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.project1.SameAs(test.project2)
			want := test.want

			assert.Equal(t, want, got)
		})
	}
}

func TestProject_WebURL(t *testing.T) {
	project := Project{
		Name:     "foo",
		Owner:    "bar",
		Host:     "github.com",
		Protocol: "https",
	}

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

func TestProject_GitURL(t *testing.T) {
	os.Setenv("HUB_PROTOCOL", "https")
	defer os.Setenv("HUB_PROTOCOL", "")

	project := Project{
		Name:  "foo",
		Owner: "bar",
		Host:  "github.com",
	}

	url := project.GitURL("gh", "jingweno", false)
	assert.Equal(t, "https://github.com/jingweno/gh.git", url)

	os.Setenv("HUB_PROTOCOL", "git")
	url = project.GitURL("gh", "jingweno", false)
	assert.Equal(t, "git://github.com/jingweno/gh.git", url)

	os.Setenv("HUB_PROTOCOL", "ssh")
	url = project.GitURL("gh", "jingweno", true)
	assert.Equal(t, "git@github.com:jingweno/gh.git", url)

	url = project.GitURL("gh", "jingweno", true)
	assert.Equal(t, "git@github.com:jingweno/gh.git", url)
}

func TestProject_GitURLEnterprise(t *testing.T) {
	project := Project{
		Name:  "foo",
		Owner: "bar",
		Host:  "https://github.corporate.com",
	}

	defer os.Setenv("HUB_PROTOCOL", "")

	os.Setenv("HUB_PROTOCOL", "https")
	url := project.GitURL("gh", "jingweno", false)
	assert.Equal(t, "https://github.corporate.com/jingweno/gh.git", url)

	os.Setenv("HUB_PROTOCOL", "ssh")
	url = project.GitURL("gh", "jingweno", false)
	assert.Equal(t, "git@github.corporate.com:jingweno/gh.git", url)

	os.Setenv("HUB_PROTOCOL", "git")
	url = project.GitURL("gh", "jingweno", false)
	assert.Equal(t, "git://github.corporate.com/jingweno/gh.git", url)

	url = project.GitURL("gh", "jingweno", true)
	assert.Equal(t, "git@github.corporate.com:jingweno/gh.git", url)
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
