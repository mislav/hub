package sawyer

import (
	"net/url"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/jingweno/go-sawyer/hypermedia"
)

var endpoints = map[string]map[string]string{
	"http://api.github.com": map[string]string{
		"user":                "http://api.github.com/user",
		"/user":               "http://api.github.com/user",
		"http://api.com/user": "http://api.com/user",
	},
	"http://api.github.com/api/v1": map[string]string{
		"user":                "http://api.github.com/api/v1/user",
		"/user":               "http://api.github.com/user",
		"http://api.com/user": "http://api.com/user",
	},
}

func TestResolve(t *testing.T) {
	for endpoint, tests := range endpoints {
		client, err := NewFromString(endpoint, nil)
		if err != nil {
			t.Fatal(err.Error())
		}

		for relative, result := range tests {
			u, err := url.Parse(relative)
			if err != nil {
				t.Error(err.Error())
				break
			}

			abs := client.ResolveReference(u)
			if absurl := abs.String(); result != absurl {
				t.Errorf("Bad absolute URL %s for %s + %s == %s", absurl, endpoint, relative, result)
			}
		}
	}
}

func TestResolveWithNoHeader(t *testing.T) {
	client, err := NewFromString("http://api.github.com", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	req, _ := client.NewRequest("")
	assert.Equal(t, 0, len(req.Header))

	req.Header.Set("Cache-Control", "private")
	assert.Equal(t, 1, len(req.Header))
	assert.Equal(t, 0, len(client.Header))
}

func TestResolveWithHeader(t *testing.T) {
	client, err := NewFromString("http://api.github.com", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	client.Header.Set("Cache-Control", "private")

	req, _ := client.NewRequest("")
	assert.Equal(t, 1, len(req.Header))
	assert.Equal(t, "private", req.Header.Get("Cache-Control"))
}

func TestResolveClientQuery(t *testing.T) {
	client, err := NewFromString("http://api.github.com", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	u, err := client.ResolveReferenceString("/foo?a=1")
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, "http://api.github.com/foo?a=1", u)
}

func TestResolveClientQueryWithClientQuery(t *testing.T) {
	client, err := NewFromString("http://api.github.com?a=1&b=1", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, "1", client.Query.Get("a"))
	assert.Equal(t, "1", client.Query.Get("b"))

	client.Query.Set("b", "2")
	client.Query.Set("c", "3")
	u, err := client.ResolveReferenceString("/foo?d=4")
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, "http://api.github.com/foo?a=1&b=2&c=3&d=4", u)
}

func TestResolveClientRelativeReference(t *testing.T) {
	client, err := NewFromString("http://github.enterprise.com/api/v3/", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	u, err := client.ResolveReferenceString("users")
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, "http://github.enterprise.com/api/v3/users", u)
}

func TestResolveClientRelativeHyperlink(t *testing.T) {
	client, err := NewFromString("http://github.enterprise.com/api/v3/", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	link := hypermedia.Hyperlink("repos/{repo}")
	expanded, err := link.Expand(hypermedia.M{"repo": "foo"})

	u, err := client.ResolveReferenceString(expanded.String())
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, "http://github.enterprise.com/api/v3/repos/foo", u)
}
