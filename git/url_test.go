package git

import (
	"testing"

	"github.com/bmizerany/assert"
)

func createURLParser() *URLParser {
	c := make(SSHConfig)
	c["github.com"] = "ssh.github.com"
	c["gh"] = "github.com"
	c["git.company.com"] = "ssh.git.company.com"

	return &URLParser{c}
}

func TestURLParser_ParseURL_HTTPURL(t *testing.T) {
	p := createURLParser()

	u, err := p.Parse("https://github.com/octokit/go-octokit.git")
	assert.Equal(t, nil, err)
	assert.Equal(t, "github.com", u.Host)
	assert.Equal(t, "https", u.Scheme)
	assert.Equal(t, "/octokit/go-octokit.git", u.Path)
}

func TestURLParser_ParseURL_GitURL(t *testing.T) {
	p := createURLParser()

	u, err := p.Parse("git://github.com/octokit/go-octokit.git")
	assert.Equal(t, nil, err)
	assert.Equal(t, "github.com", u.Host)
	assert.Equal(t, "git", u.Scheme)
	assert.Equal(t, "/octokit/go-octokit.git", u.Path)

	u, err = p.Parse("https://git.company.com/octokit/go-octokit.git")
	assert.Equal(t, nil, err)
	assert.Equal(t, "git.company.com", u.Host)
	assert.Equal(t, "https", u.Scheme)
	assert.Equal(t, "/octokit/go-octokit.git", u.Path)

	u, err = p.Parse("git://git.company.com/octokit/go-octokit.git")
	assert.Equal(t, nil, err)
	assert.Equal(t, "git.company.com", u.Host)
	assert.Equal(t, "git", u.Scheme)
	assert.Equal(t, "/octokit/go-octokit.git", u.Path)
}

func TestURLParser_ParseURL_SSHURL(t *testing.T) {
	p := createURLParser()

	u, err := p.Parse("git@github.com:lostisland/go-sawyer.git")
	assert.Equal(t, nil, err)
	assert.Equal(t, "github.com", u.Host)
	assert.Equal(t, "ssh", u.Scheme)
	assert.Equal(t, "git", u.User.Username())
	assert.Equal(t, "/lostisland/go-sawyer.git", u.Path)

	u, err = p.Parse("gh:octokit/go-octokit")
	assert.Equal(t, nil, err)
	assert.Equal(t, "github.com", u.Host)
	assert.Equal(t, "ssh", u.Scheme)
	assert.Equal(t, "/octokit/go-octokit", u.Path)

	u, err = p.Parse("git@git.company.com:octokit/go-octokit")
	assert.Equal(t, nil, err)
	assert.Equal(t, "ssh.git.company.com", u.Host)
	assert.Equal(t, "ssh", u.Scheme)
	assert.Equal(t, "/octokit/go-octokit", u.Path)
}

func TestURLParser_ParseURL_LocalPath(t *testing.T) {
	p := createURLParser()

	u, err := p.Parse("/path/to/repo.git")
	assert.Equal(t, nil, err)
	assert.Equal(t, "", u.Host)
	assert.Equal(t, "", u.Scheme)
	assert.Equal(t, "/path/to/repo.git", u.Path)

	u, err = p.Parse(`c:\path\to\repo.git`)
	assert.Equal(t, nil, err)
	assert.Equal(t, `c:\path\to\repo.git`, u.String())
}
