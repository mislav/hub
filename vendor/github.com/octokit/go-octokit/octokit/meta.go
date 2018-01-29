package octokit

import (
	"encoding/json"
	"net"

	"github.com/jingweno/go-sawyer/hypermedia"
)

// https://developer.github.com/v3/meta/
var (
	MetaURL = Hyperlink("/meta")
)

// Meta return an APIInfo with the current API meta information
//
// https://developer.github.com/v3/meta/#meta
func (c *Client) Meta(uri *Hyperlink) (info APIInfo, result *Result) {
	url, err := uri.Expand(nil)
	if err != nil {
		return info, &Result{Err: err}
	}
	var meta meta
	result = c.get(url, &meta)
	if !result.HasError() {
		info = meta.transform()
	}
	return
}

type ipNets []*net.IPNet

func (i *ipNets) UnmarshalJSON(raw []byte) error {
	*i = (*i)[:0]
	var ss []string
	if err := json.Unmarshal(raw, &ss); err != nil {
		return err
	}
	for _, s := range ss {
		_, ipNet, err := net.ParseCIDR(s)
		if err != nil {
			return err
		}
		*i = append(*i, ipNet)
	}
	return nil
}

type ips []net.IP

func (i *ips) UnmarshalJSON(raw []byte) error {
	*i = (*i)[:0]
	var ss []string
	if err := json.Unmarshal(raw, &ss); err != nil {
		return err
	}
	for _, s := range ss {
		*i = append(*i, net.ParseIP(s))
	}
	return nil
}

type meta struct {
	*hypermedia.HALResource

	VerifiablePasswordAuthentication bool   `json:"verifiable_password_authentication,omitempty"`
	GithubServicesSha                string `json:"github_services_sha,omitempty"`
	Hooks                            ipNets `json:"hooks,omitempty"`
	Git                              ipNets `json:"git,omitempty"`
	Pages                            ipNets `json:"pages,omitempty"`
	Importer                         ips    `json:"importer,omitempty"`
}

func (m meta) transform() (info APIInfo) {
	info.VerifiablePasswordAuthentication = m.VerifiablePasswordAuthentication
	info.GithubServicesSha = m.GithubServicesSha

	info.Hooks = ([]*net.IPNet)(m.Hooks)
	info.Git = ([]*net.IPNet)(m.Git)
	info.Pages = ([]*net.IPNet)(m.Pages)
	info.Importer = ([]net.IP)(m.Importer)

	return
}

// APIInfo contains the information described in https://developer.github.com/v3/meta/#body
type APIInfo struct {
	*hypermedia.HALResource

	VerifiablePasswordAuthentication bool         `json:"verifiable_password_authentication,omitempty"`
	GithubServicesSha                string       `json:"github_services_sha,omitempty"`
	Hooks                            []*net.IPNet `json:"hooks,omitempty"`
	Git                              []*net.IPNet `json:"git,omitempty"`
	Pages                            []*net.IPNet `json:"pages,omitempty"`
	Importer                         []net.IP     `json:"importer,omitempty"`
}
