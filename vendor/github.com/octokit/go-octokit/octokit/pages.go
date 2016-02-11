package octokit

import (
	"github.com/jingweno/go-sawyer/hypermedia"
	"time"
)

// URL templates for actions taken on the pages in repositories
//
// https://developer.github.com/v3/repos/pages/
var (
	PagesURL            = Hyperlink("/repos/{owner}/{repo}/pages")
	PagesBuildsURL      = Hyperlink("/repos/{owner}/{repo}/pages/builds")
	PagesLatestBuildURL = Hyperlink("/repos/{owner}/{repo}/pages/builds/latest")
)

// Pages creates a PagesService to access page information including various
// versions of build
//
// https://developer.github.com/v3/repos/pages/
func (c *Client) Pages() *PagesService {
	return &PagesService{client: c}
}

// A service to return page information including various versions of build
type PagesService struct {
	client *Client
}

// PageInfo gets the information about a Pages site
//
// https://developer.github.com/v3/repos/pages/#get-information-about-a-pages-site
func (g *PagesService) PageInfo(uri *Hyperlink, uriParams M) (page *PageInfo,
	result *Result) {
	url, err := ExpandWithDefault(uri, &PagesURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}
	result = g.client.get(url, &page)
	return
}

// PageBuilds lists the builds of a given page
// https://developer.github.com/v3/repos/pages/#list-pages-builds
func (g *PagesService) PageBuilds(uri *Hyperlink, uriParams M) (builds []PageBuild,
	result *Result) {
	url, err := ExpandWithDefault(uri, &PagesBuildsURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}
	result = g.client.get(url, &builds)
	return
}

// PageBuildLatest gets the latest build for a page
//
// https://developer.github.com/v3/repos/pages/#list-latest-pages-build
func (g *PagesService) PageBuildLatest(uri *Hyperlink, uriParams M) (build *PageBuild,
	result *Result) {
	url, err := ExpandWithDefault(uri, &PagesLatestBuildURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}
	result = g.client.get(url, &build)
	return
}

type PageInfo struct {
	*hypermedia.HALResource

	URL       string `json:"url,omitempty"`
	Status    string `json:"status,omitempty"`
	Cname     string `json:"cname,omitempty"`
	Custom404 bool   `json:"custom_404,omitempty"`
}

type PageBuild struct {
	*hypermedia.HALResource

	URL       string       `json:"url,omitempty"`
	Status    string       `json:"status,omitempty"`
	Error     *ErrorObject `json:"error,omitempty"`
	Pusher    *User        `json:"pusher,omitempty"`
	Commit    string       `json:"commit,omitempty"`
	Duration  int          `json:"duration,omitempty"`
	CreatedAt *time.Time   `json:"created_at,omitempty"`
	UpdatedAt *time.Time   `json:"updated_at,omitempty"`
}
