package octokit

import (
	"net/url"

	"github.com/jingweno/go-sawyer/hypermedia"
)

// EmailURL is an address for accessing email addresses for the current user
//
// https://developer.github.com/v3/users/emails/
var EmailUrl = Hyperlink("user/emails")

// Create a EmailsService with the base url.URL
//
// https://developer.github.com/v3/users/emails/
func (c *Client) Emails(url *url.URL) (emails *EmailsService) {
	emails = &EmailsService{client: c, URL: url}
	return
}

// A service to return user email addresses
type EmailsService struct {
	client *Client
	URL    *url.URL
}

// Get a list of email addresses for the current user
//
// https://developer.github.com/v3/users/emails/#list-email-addresses-for-a-user
func (e *EmailsService) All() (emails []Email, result *Result) {
	result = e.client.get(e.URL, &emails)
	return
}

// Adds a list of email address(es) for the current user
//
// https://developer.github.com/v3/users/emails/#add-email-addresses
func (e *EmailsService) Create(params interface{}) (emails []Email, result *Result) {
	result = e.client.post(e.URL, params, &emails)
	return
}

// Deletes email address(es) for the current user
//
// https://developer.github.com/v3/users/emails/#delete-email-addresses
func (e *EmailsService) Delete(params interface{}) (result *Result) {
	result = e.client.delete(e.URL, params, nil)
	return
}

type Email struct {
	*hypermedia.HALResource

	Email    string `json:"email,omitempty"`
	Verified bool   `json:"verified,omitempty"`
	Primary  bool   `json:"primary,omitempty"`
}
