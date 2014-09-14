package octokit

import (
	"net/url"
	"time"

	"github.com/jingweno/go-sawyer/hypermedia"
)

var (
	CurrentUserURL = Hyperlink("user")
	UserURL        = Hyperlink("users{/user}")
)

// Create a UsersService with the base url.URL
func (c *Client) Users(url *url.URL) (users *UsersService) {
	users = &UsersService{client: c, URL: url}
	return
}

// A service to return user records
type UsersService struct {
	client *Client
	URL    *url.URL
}

// Get a user based on UserService#URL
func (u *UsersService) One() (user *User, result *Result) {
	result = u.client.get(u.URL, &user)
	return
}

// Update a user based on UserService#URL
func (u *UsersService) Update(params interface{}) (user *User, result *Result) {
	result = u.client.put(u.URL, params, &user)
	return
}

// Get a list of users based on UserService#URL
func (u *UsersService) All() (users []User, result *Result) {
	result = u.client.get(u.URL, &users)
	return
}

type User struct {
	*hypermedia.HALResource

	SiteAdmin         bool       `json:"site_admin,omitempty"`
	Login             string     `json:"login,omitempty"`
	ID                int        `json:"id,omitempty"`
	AvatarURL         string     `json:"avatar_url,omitempty"`
	GravatarID        string     `json:"gravatar_id,omitempty"`
	URL               string     `json:"url,omitempty"`
	Name              string     `json:"name,omitempty"`
	Company           string     `json:"company,omitempty"`
	Blog              string     `json:"blog,omitempty"`
	Location          string     `json:"location,omitempty"`
	Email             string     `json:"email,omitempty"`
	Hireable          bool       `json:"hireable,omitempty"`
	Bio               string     `json:"bio,omitempty"`
	PublicRepos       int        `json:"public_repos,omitempty"`
	Followers         int        `json:"followers,omitempty"`
	Following         int        `json:"following,omitempty"`
	HTMLURL           string     `json:"html_url,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	Type              string     `json:"type,omitempty"`
	FollowingURL      Hyperlink  `json:"following_url,omitempty"`
	FollowersURL      Hyperlink  `json:"followers_url,omitempty"`
	GistsURL          Hyperlink  `json:"gists_url,omitempty"`
	StarredURL        Hyperlink  `json:"starred_url,omitempty"`
	SubscriptionsURL  Hyperlink  `json:"subscriptions_url,omitempty"`
	OrganizationsURL  Hyperlink  `json:"organizations_url,omitempty"`
	ReposURL          Hyperlink  `json:"repos_url,omitempty"`
	EventsURL         Hyperlink  `json:"events_url,omitempty"`
	ReceivedEventsURL Hyperlink  `json:"received_events_url,omitempty"`
}
