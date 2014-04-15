package github

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestClient_ApiEndpoint(t *testing.T) {
	gh := &Client{Credentials: &Credentials{Host: "github.com"}}
	assert.Equal(t, "https://api.github.com", gh.apiEndpoint())

	gh = &Client{Credentials: &Credentials{Host: "github.corporate.com"}}
	assert.Equal(t, "https://github.corporate.com", gh.apiEndpoint())

	gh = &Client{Credentials: &Credentials{Host: "http://github.corporate.com"}}
	assert.Equal(t, "http://github.corporate.com", gh.apiEndpoint())
}
