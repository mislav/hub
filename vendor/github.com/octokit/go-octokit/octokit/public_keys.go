package octokit

import (
	"time"

	"github.com/jingweno/go-sawyer/hypermedia"
)

var (
	CurrentPublicKeyUrl = Hyperlink("/user/keys{/id}")
	PublicKeyUrl        = Hyperlink("/users/{user}/keys")
)

// Create a PublicKeysService
func (c *Client) PublicKeys() (k *PublicKeysService) {
	k = &PublicKeysService{client: c}
	return
}

// A service to return user public keys
type PublicKeysService struct {
	client *Client
}

// Get a list of keys for the user
//
// https://developer.github.com/v3/repos/keys/#list-deploy-keys
func (k *PublicKeysService) All(uri *Hyperlink, uriParams M) (keys []Key, result *Result) {
	url, err := ExpandWithDefault(uri, &CurrentPublicKeyUrl, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = k.client.get(url, &keys)
	return
}

// Get a the data for one key for the current user
//
// https://developer.github.com/v3/repos/keys/#get-a-deploy-key
func (k *PublicKeysService) One(uri *Hyperlink, uriParams M) (key *Key, result *Result) {
	url, err := ExpandWithDefault(uri, &CurrentPublicKeyUrl, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = k.client.get(url, &key)
	return
}

// Creates a new public key for the current user
//
// https://developer.github.com/v3/repos/keys/#add-a-new-deploy-key
func (k *PublicKeysService) Create(uri *Hyperlink, uriParams M, requestParams interface{}) (key *Key, result *Result) {
	url, err := ExpandWithDefault(uri, &CurrentPublicKeyUrl, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = k.client.post(url, requestParams, &key)
	return
}

// Removes a public key for the current user
//
// https://developer.github.com/v3/repos/keys/#remove-a-deploy-key
func (k *PublicKeysService) Delete(uri *Hyperlink, uriParams M) (success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &CurrentPublicKeyUrl, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}

	result = k.client.delete(url, nil, nil)
	success = (result.Response.StatusCode == 204)
	return
}

type Key struct {
	*hypermedia.HALResource

	Id        int        `json:"id,omitempty"`
	Key       string     `json:"key,omitempty"`
	URL       string     `json:"url,omitempty"`
	Title     string     `json:"title,omitempty"`
	Verified  bool       `json:"verified,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
