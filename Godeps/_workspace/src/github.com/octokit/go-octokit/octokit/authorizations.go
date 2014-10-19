package octokit

import (
	"net/url"
	"time"

	"github.com/jingweno/go-sawyer/hypermedia"
)

var (
	AuthorizationsURL = Hyperlink("authorizations{/id}")
)

// Create a AuthorizationsService with the base url.URL
func (c *Client) Authorizations(url *url.URL) (auths *AuthorizationsService) {
	auths = &AuthorizationsService{client: c, URL: url}
	return
}

type AuthorizationsService struct {
	client *Client
	URL    *url.URL
}

func (a *AuthorizationsService) One() (auth *Authorization, result *Result) {
	result = a.client.get(a.URL, &auth)
	return
}

func (a *AuthorizationsService) All() (auths []Authorization, result *Result) {
	result = a.client.get(a.URL, &auths)
	return
}

func (a *AuthorizationsService) Create(params interface{}) (auth *Authorization, result *Result) {
	result = a.client.post(a.URL, params, &auth)
	return
}

type Authorization struct {
	*hypermedia.HALResource

	ID        int       `json:"id,omitempty"`
	URL       string    `json:"url,omitempty"`
	App       App       `json:"app,omitempty"`
	Token     string    `json:"token,omitempty"`
	Note      string    `json:"note,omitempty"`
	NoteURL   string    `json:"note_url,omitempty"`
	Scopes    []string  `json:"scopes,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type App struct {
	*hypermedia.HALResource

	ClientID string `json:"client_id,omitempty"`
	URL      string `json:"url,omitempty"`
	Name     string `json:"name,omitempty"`
}

type AuthorizationParams struct {
	Scopes       []string `json:"scopes,omitempty"`
	Note         string   `json:"note,omitempty"`
	NoteURL      string   `json:"note_url,omitempty"`
	ClientID     string   `json:"client_id,omitempty"`
	ClientSecret string   `json:"client_secret,omitempty"`
}
