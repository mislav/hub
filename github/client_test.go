package github

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/github/hub/fixtures"
	"github.com/octokit/go-octokit/octokit"
)

func TestClient_newOctokitClient(t *testing.T) {
	c := NewClient("github.com")
	cc := c.newOctokitClient(nil)
	assert.Equal(t, "https://api.github.com/", cc.Endpoint.String())

	c = NewClient("github.corporate.com")
	cc = c.newOctokitClient(nil)
	assert.Equal(t, "https://github.corporate.com/", cc.Endpoint.String())
}

func TestClient_newTLSConfig(t *testing.T) {
	c := NewClient("github.corporate.com")
	tlsConfig, err := c.newTLSConfig()
	assert.Equal(t, nil, err)

	// We don't have a client certificate configured, so config should come
	// back as nil.
	assert.Equal(t, (*tls.Config)(nil), tlsConfig)

	c.Host.ClientCertificate = fixtures.Path("cert", "client.p12")
	tlsConfig, err = c.newTLSConfig()
	assert.Equal(t, nil, err)
}

func TestClient_FormatError(t *testing.T) {
	e := &octokit.ResponseError{
		Response: &http.Response{
			StatusCode: 401,
			Status:     "401 Not Found",
		},
	}

	err := FormatError("action", e)
	assert.Equal(t, "Error action: Not Found (HTTP 401)", fmt.Sprintf("%s", err))

	e = &octokit.ResponseError{
		Response: &http.Response{
			StatusCode: 422,
			Status:     "422 Unprocessable Entity",
		},
		Message: "error message",
	}
	err = FormatError("action", e)
	assert.Equal(t, "Error action: Unprocessable Entity (HTTP 422)\nerror message", fmt.Sprintf("%s", err))
}

func TestAuthTokenNote(t *testing.T) {
	note, err := authTokenNote(1)
	assert.Equal(t, nil, err)

	reg := regexp.MustCompile("hub for (.+)@(.+)")
	assert.T(t, reg.MatchString(note))

	note, err = authTokenNote(2)
	assert.Equal(t, nil, err)

	reg = regexp.MustCompile("hub for (.+)@(.+) 2")
	assert.T(t, reg.MatchString(note))

}
