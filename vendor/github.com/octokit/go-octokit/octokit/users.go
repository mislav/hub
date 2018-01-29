package octokit

import (
	"net/url"
	"time"

	"github.com/jingweno/go-sawyer/hypermedia"
)

//Hyperlinks to various ways of accessing users on github.
//CurrentUserURL is the address for the current user.
//UserURL is a template for the address any particular user or all users.
//
// https://developer.github.com/v3/users/
var (
	CurrentUserURL = Hyperlink("user")
	UserURL        = Hyperlink("users{/user}")
)

// Users creates a UsersService with a base url
//
// https://developer.github.com/v3/users/
func (c *Client) Users(url *url.URL) (users *UsersService) {
	users = &UsersService{client: c, URL: url}
	return
}

// UsersService is a service providing access to user records from a particular url
type UsersService struct {
	client *Client
	URL    *url.URL
}

// One gets a specific user record based on the url of the service
//
// https://developer.github.com/v3/users/#get-a-single-user
// https://developer.github.com/v3/users/#get-the-authenticated-user
func (u *UsersService) One() (user *User, result *Result) {
	result = u.client.get(u.URL, &user)
	return
}

// Update modifies a user record specified in the User struct as parameters on the
// service url
//
// https://developer.github.com/v3/users/#update-the-authenticated-user
func (u *UsersService) Update(params interface{}) (user *User, result *Result) {
	result = u.client.put(u.URL, params, &user)
	return
}

// All gets a list of all user records associated with the url of the service
//
// https://developer.github.com/v3/users/#get-all-users
func (u *UsersService) All() (users []User, result *Result) {
	result = u.client.get(u.URL, &users)
	return
}

// User represents the full user record of a particular user on GitHub
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
