package octokit

import (
	"time"
)

type User struct {
	Login       string    `json:"login"`
	Id          int       `json:"id"`
	AvatarUrl   string    `json:"avatar_url"`
	GravatarId  string    `json:"gravatar_id"`
	Url         string    `json:"url"`
	Name        string    `json:"name"`
	Company     string    `json:"company"`
	Blog        string    `json:"blog"`
	Location    string    `json:"location"`
	Email       string    `json:"email"`
	Hireable    bool      `json:"hireable"`
	Bio         string    `json:"bio"`
	PublicRepos int       `json:"public_repos"`
	PublicGists int       `json:"jsonpublic_gists"`
	Followers   int       `json:"followers"`
	Following   int       `json:"following"`
	HtmlUrl     string    `json:"html_url"`
	CreatedAt   time.Time `json:"created_at"`
	Type        string    `json:"jsontype"`
}

func (c *Client) AuthenticatedUser() (*User, error) {
	body, err := c.get("user", nil)
	if err != nil {
		return nil, err
	}

	var authUser User
	err = jsonUnmarshal(body, &authUser)

	return &authUser, err
}
