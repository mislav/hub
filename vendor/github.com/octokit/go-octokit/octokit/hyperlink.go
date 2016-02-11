package octokit

import (
	"net/url"

	"github.com/jingweno/go-sawyer/hypermedia"
)

// M represents a map of values to expand a Hyperlink. The keys in M are elements
// of the template to be replaced with the associated value in the map.
type M map[string]interface{}

// Hyperlink is a string url.  If it is a uri template, it can be converted to
// a full URL with Expand().
type Hyperlink string

// Expand utilizes the sawyer expand method to convert a URI template into a full
// URL
func (l Hyperlink) Expand(m M) (u *url.URL, err error) {
	sawyerHyperlink := hypermedia.Hyperlink(string(l))
	u, err = sawyerHyperlink.Expand(hypermedia.M(m))
	return
}

// Expands a link with possible, otherwise it expands the default link
func ExpandWithDefault(link *Hyperlink, defaultLink *Hyperlink, params M) (u *url.URL, err error) {
	if link == nil {
		link = defaultLink
	}

	return link.Expand(params)
}
