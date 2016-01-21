package hypermedia

import (
	"bytes"
	"encoding/json"
	"github.com/bmizerany/assert"
	"testing"
)

func TestReflectRelations(t *testing.T) {
	input := `
{ "Login": "bob"
, "Url": "/self"
, "FooUrl": "/foo"
, "FooBarUrl": "/bar"
, "whatever": "/whatevs"
, "HomepageUrl": "http://example.com"
}`

	user := &ReflectedUser{}
	decode(t, input, user)

	rels := HyperFieldDecoder(user)
	assert.Equal(t, 4, len(rels))
	assert.Equal(t, "/self", string(rels["Url"]))
	assert.Equal(t, "/foo", string(rels["FooUrl"]))
	assert.Equal(t, "/bar", string(rels["FooBarUrl"]))
	assert.Equal(t, "/whatevs", string(rels["whatevs"]))

	rel, err := rels.Rel("FooUrl", nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, "/foo", rel.Path)
}

func TestHALRelations(t *testing.T) {
	input := `
{ "Login": "bob"
, "Url": "/foo/bar{/arg}"
, "_links":
	{ "self": { "href": "/self" }
	, "foo": { "href": "/foo" }
	, "bar": { "href": "/bar" }
	}
}`

	user := &HypermediaUser{}
	decode(t, input, user)

	rels := HypermediaDecoder(user)
	assert.Equal(t, 3, len(rels))
	assert.Equal(t, "/self", string(rels["self"]))
	assert.Equal(t, "/foo", string(rels["foo"]))
	assert.Equal(t, "/bar", string(rels["bar"]))

	rel, err := rels.Rel("foo", nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, "/foo", rel.Path)
}

func TestExpandAbsoluteUrls(t *testing.T) {
	link := Hyperlink("/foo/bar{/arg}")
	u, err := link.Expand(M{"arg": "baz", "foo": "bar"})
	assert.Equal(t, nil, err)
	assert.Equal(t, "/foo/bar/baz", u.String())
}

func TestExpandRelativePaths(t *testing.T) {
	link := Hyperlink("foo/bar{/arg}")
	u, err := link.Expand(M{"arg": "baz", "foo": "bar"})
	assert.Equal(t, nil, err)
	assert.Equal(t, "foo/bar/baz", u.String())
}

func TestExpandNil(t *testing.T) {
	link := Hyperlink("/foo/bar{/arg}")
	u, err := link.Expand(nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, "/foo/bar", u.String())
}

func TestDecode(t *testing.T) {
	input := `
{ "Login": "bob"
, "Url": "/foo/bar{/arg}"
, "_links":
  { "self": { "href": "/foo/bar{/arg}" }
  }
}`

	user := &HypermediaUser{}
	decode(t, input, user)

	assert.Equal(t, "bob", user.Login)
	assert.Equal(t, 1, len(user.Links))

	hl := user.Url
	url, err := hl.Expand(M{"arg": "baz"})
	if err != nil {
		t.Errorf("Errors parsing %s: %s", hl, err)
	}

	assert.Equal(t, "/foo/bar/baz", url.String())

	hl = user.Links["self"].Href
	url, err = hl.Expand(M{"arg": "baz"})
	if err != nil {
		t.Errorf("Errors parsing %s: %s", hl, err)
	}
	assert.Equal(t, "/foo/bar/baz", url.String())
}

func decode(t *testing.T, input string, resource interface{}) {
	dec := json.NewDecoder(bytes.NewBufferString(input))
	err := dec.Decode(resource)
	if err != nil {
		t.Fatalf("Errors decoding json: %s", err)
	}
}

type HypermediaUser struct {
	Login string
	Url   Hyperlink
	*HALResource
}

type ReflectedUser struct {
	Login       string
	Url         Hyperlink
	FooUrl      Hyperlink
	FooBarUrl   Hyperlink
	Whatever    Hyperlink `json:"whatever" rel:"whatevs"`
	HomepageUrl string
}
