package octokit

import (
	"net/url"
	"time"

	"github.com/jingweno/go-sawyer/hypermedia"
)

// AuthorizationsURL is a template for accessing authorizations possibly associated with a
// given identification number that can be expanded to a full address.
//
// https://developer.github.com/v3/oauth_authorizations/
var AuthorizationsURL = Hyperlink("authorizations{/id}")

// Authorizations creates a AuthorizationsService with a base url.
//
// https://developer.github.com/v3/oauth_authorizations/
func (c *Client) Authorizations(url *url.URL) (auths *AuthorizationsService) {
	auths = &AuthorizationsService{client: c, URL: url}
	return
}

// AuthorizationsService is a service providing access to OAuth authorizations through
// the authorization API
type AuthorizationsService struct {
	client *Client
	URL    *url.URL
}

// One gets a specific authorization based on the url of the service.
//
// https://developer.github.com/v3/oauth_authorizations/#get-a-single-authorization
func (a *AuthorizationsService) One() (auth *Authorization, result *Result) {
	result = a.client.get(a.URL, &auth)
	return
}

// All gets a list of all authorizations associated with the url of the service.
//
// https://developer.github.com/v3/oauth_authorizations/#list-your-authorizations
func (a *AuthorizationsService) All() (auths []Authorization, result *Result) {
	result = a.client.get(a.URL, &auths)
	return
}

// Create posts a new authorization to the authorizations service url.
//
// https://developer.github.com/v3/oauth_authorizations/#create-a-new-authorization
func (a *AuthorizationsService) Create(params interface{}) (auth *Authorization, result *Result) {
	result = a.client.post(a.URL, params, &auth)
	return
}

// Authorization is a representation of an OAuth passed to or from the authorizations API
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

// App is the unit holding the associated authorization
type App struct {
	*hypermedia.HALResource

	ClientID string `json:"client_id,omitempty"`
	URL      string `json:"url,omitempty"`
	Name     string `json:"name,omitempty"`
}

// AuthorizationParams is the set of parameters used when creating a new
// authorization to be posted to the API
type AuthorizationParams struct {
	Scopes       []string `json:"scopes,omitempty"`
	Note         string   `json:"note,omitempty"`
	NoteURL      string   `json:"note_url,omitempty"`
	ClientID     string   `json:"client_id,omitempty"`
	ClientSecret string   `json:"client_secret,omitempty"`
}
