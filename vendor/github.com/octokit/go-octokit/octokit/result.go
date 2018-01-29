package octokit

import (
	"strconv"
	"strings"
	"time"

	"github.com/jingweno/go-sawyer/mediaheader"
)

const (
	oauthScopes         = "X-OAuth-Scopes"
	oauthAcceptedScopes = "X-OAuth-Accepted-Scopes"
	rateLimitRemaining  = "X-RateLimit-Remaining"
	rateLimitReset      = "X-RateLimit-Reset"
)

type pageable struct {
	NextPage  *Hyperlink
	LastPage  *Hyperlink
	FirstPage *Hyperlink
	PrevPage  *Hyperlink
}

// Result is a pageable set of data, with hyperlinks to the first, last,
// previous, and next pages, containing a response to some request and
// associated error, if any
type Result struct {
	Response *Response
	Err      error
	pageable
}

// HasError returns true if the error field of the Result is not nil; false
// otherwise
func (r *Result) HasError() bool {
	return r.Err != nil
}

// Error returns the string representation of the error if it exists; the
// empty string is returned otherwise
func (r *Result) Error() string {
	if r.Err != nil {
		return r.Err.Error()
	}

	return ""
}

func (r *Result) RateLimitReset() *time.Time {
	epoc := r.Response.Header.Get(rateLimitReset)
	if epoc == "" {
		return nil
	}

	reset, err := strconv.ParseInt(epoc, 10, 64)
	if err != nil {
		return nil
	}

	t := time.Unix(reset, 0)
	return &t
}

func (r *Result) RateLimitRemaining() int {
	rate, err := strconv.Atoi(r.Response.Header.Get(rateLimitRemaining))
	if err != nil {
		rate = defaultRateLimit(r.Response)
	}
	return rate
}

func (r *Result) RawScopes() string {
	return r.Response.Header.Get(oauthScopes)
}

func (r *Result) Scopes() []string {
	return strings.Split(r.RawScopes(), ", ")
}

func (r *Result) RawAcceptedScopes() string {
	return r.Response.Header.Get(oauthAcceptedScopes)
}

func (r *Result) AcceptedScopes() []string {
	return strings.Split(r.RawAcceptedScopes(), ", ")
}

func (r *Result) ValidScope(scope string) bool {
	for _, s := range r.Scopes() {
		if s == scope {
			return true
		}
	}
	return false
}

func newResult(resp *Response, err error) *Result {
	pageable := pageable{}
	if resp != nil {
		fillPageable(&pageable, resp.MediaHeader)
	}

	return &Result{Response: resp, pageable: pageable, Err: err}
}

func fillPageable(pageable *pageable, header *mediaheader.MediaHeader) {
	if link, ok := header.Relations["next"]; ok {
		l := Hyperlink(link)
		pageable.NextPage = &l
	}

	if link, ok := header.Relations["prev"]; ok {
		l := Hyperlink(link)
		pageable.PrevPage = &l
	}

	if link, ok := header.Relations["first"]; ok {
		l := Hyperlink(link)
		pageable.FirstPage = &l
	}

	if link, ok := header.Relations["last"]; ok {
		l := Hyperlink(link)
		pageable.LastPage = &l
	}
}

func defaultRateLimit(r *Response) int {
	if r.Request != nil {
		h := r.Request.URL.Host
		if !strings.HasSuffix(gitHubAPIURL, h) {
			return -1
		}
	}
	return 60
}
