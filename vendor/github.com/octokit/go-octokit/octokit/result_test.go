package octokit

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jingweno/go-sawyer/hypermedia"
	"github.com/jingweno/go-sawyer/mediaheader"
	"github.com/stretchr/testify/assert"
)

func TestNewResult_Pageable(t *testing.T) {
	resp := newTestResponse()
	result := newResult(resp, nil)

	assert.Equal(t, "/path", string(*result.NextPage))
	assert.Nil(t, result.PrevPage)
	assert.Nil(t, result.LastPage)
	assert.Nil(t, result.FirstPage)
}

func TestRateLimitReset(t *testing.T) {
	etime := time.Unix(1428697849, 0)
	cases := []struct {
		epoc     string
		expected *time.Time
	}{
		{"", nil},
		{"asdf", nil},
		{"1428697849", &etime},
	}

	for _, c := range cases {
		headers := http.Header{}
		headers.Set(rateLimitReset, c.epoc)

		resp := newTestResponse()
		resp.Header = headers

		result := newResult(resp, nil)
		assert.Equal(t, c.expected, result.RateLimitReset())
	}
}

func TestRateLimitRemaining(t *testing.T) {
	cases := []struct {
		host     string
		rate     string
		expected int
	}{
		{"api.github.com", "", 60},
		{"api.github.com", "asdf", 60},
		{"api.github.com", "3400", 3400},
		{"github.evilcorp.com", "", -1},
		{"github.evilcorp.com", "asdf", -1},
		{"github.evilcorp.com", "3400", 3400},
	}

	for _, c := range cases {
		headers := http.Header{}
		headers.Set(rateLimitRemaining, c.rate)
		req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s/user", c.host), nil)

		resp := newTestResponse()
		resp.Request = req
		resp.Header = headers

		result := newResult(resp, nil)
		assert.Equal(t, c.expected, result.RateLimitRemaining())
	}
}

func TestScopes(t *testing.T) {
	cases := []struct {
		actual   string
		expected []string
	}{
		{"user", []string{"user"}},
		{"user, repo", []string{"user", "repo"}},
	}

	for _, c := range cases {
		headers := http.Header{}
		headers.Set(oauthScopes, c.actual)

		resp := newTestResponse()
		resp.Header = headers

		result := newResult(resp, nil)
		assert.Equal(t, c.expected, result.Scopes())
		assert.Equal(t, c.actual, result.RawScopes())
	}
}

func TestAcceptedScopes(t *testing.T) {
	cases := []struct {
		actual   string
		expected []string
	}{
		{"user", []string{"user"}},
		{"user, repo", []string{"user", "repo"}},
	}

	for _, c := range cases {
		headers := http.Header{}
		headers.Set(oauthAcceptedScopes, c.actual)

		resp := newTestResponse()
		resp.Header = headers

		result := newResult(resp, nil)
		assert.Equal(t, c.expected, result.AcceptedScopes())
		assert.Equal(t, c.actual, result.RawAcceptedScopes())
	}
}

func TestValidScope(t *testing.T) {
	headers := http.Header{}
	headers.Set(oauthScopes, "user, repo")

	resp := newTestResponse()
	resp.Header = headers

	result := newResult(resp, nil)
	assert.True(t, result.ValidScope("user"))
	assert.False(t, result.ValidScope("org"))
}

func newTestResponse() *Response {
	return &Response{
		Response: &http.Response{},
		MediaHeader: &mediaheader.MediaHeader{
			Relations: hypermedia.Relations{"next": hypermedia.Hyperlink("/path")},
		},
	}
}
