package octokit

import (
	"github.com/lostisland/go-sawyer/hypermedia"
	"net/url"
)

type M map[string]interface{}

type Hyperlink string

func (l Hyperlink) Expand(m M) (u *url.URL, err error) {
	sawyerHyperlink := hypermedia.Hyperlink(string(l))
	u, err = sawyerHyperlink.Expand(hypermedia.M(m))
	return
}
