package octokit

import (
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/hypermedia"
	"net/http"
	"net/url"
	"os"
)

func NewClient(authMethod AuthMethod) *Client {
	return NewClientWith(gitHubAPIURL, nil, authMethod)
}

func NewClientWith(baseURL string, httpClient *http.Client, authMethod AuthMethod) *Client {
	client, _ := sawyer.NewFromString(baseURL, httpClient)
	return &Client{sawyerClient: client, UserAgent: userAgent, AuthMethod: authMethod}
}

type Client struct {
	UserAgent    string
	AuthMethod   AuthMethod
	sawyerClient *sawyer.Client
	rootRels     hypermedia.Relations
}

func (c *Client) NewRequest(urlStr string) (req *Request, err error) {
	sawyerReq, err := c.newSawyerRequest(urlStr)
	if err != nil {
		return
	}

	req = &Request{sawyerReq: sawyerReq}
	return
}

func (c *Client) head(url *url.URL, output interface{}) (result *Result) {
	return sendRequest(c, url, func(req *Request) (*Response, error) {
		return req.Head(output)
	})
}

func (c *Client) get(url *url.URL, output interface{}) (result *Result) {
	return sendRequest(c, url, func(req *Request) (*Response, error) {
		return req.Get(output)
	})
}

func (c *Client) post(url *url.URL, input interface{}, output interface{}) (result *Result) {
	return sendRequest(c, url, func(req *Request) (*Response, error) {
		return req.Post(input, output)
	})
}

func (c *Client) put(url *url.URL, input interface{}, output interface{}) *Result {
	return sendRequest(c, url, func(req *Request) (*Response, error) {
		return req.Put(input, output)
	})
}

func (c *Client) delete(url *url.URL, output interface{}) (result *Result) {
	return sendRequest(c, url, func(req *Request) (*Response, error) {
		return req.Delete(output)
	})
}

func (c *Client) patch(url *url.URL, input interface{}, output interface{}) (result *Result) {
	return sendRequest(c, url, func(req *Request) (*Response, error) {
		return req.Patch(input, output)
	})
}

func (c *Client) upload(uploadUrl *url.URL, asset *os.File, contentType string) (result *Result) {
	req, err := c.newSawyerRequest(uploadUrl.String())
	if err != nil {
		result = newResult(nil, err)
		return
	}

	fi, err := asset.Stat()
	if err != nil {
		result = newResult(nil, err)
		return
	}

	req.Header.Add("Content-Type", contentType)
	req.ContentLength = fi.Size()

	req.Body = asset
	sawyerResp := req.Post()

	resp, err := NewResponse(sawyerResp)
	return newResult(resp, err)
}

func (c *Client) newSawyerRequest(urlStr string) (sawyerReq *sawyer.Request, err error) {
	sawyerReq, err = c.sawyerClient.NewRequest(urlStr)
	if err != nil {
		return
	}

	sawyerReq.Header.Add("Accept", defaultMediaType)
	sawyerReq.Header.Add("User-Agent", c.UserAgent)

	c.addAuthenticationHeaders(sawyerReq.Header)

	return
}

func sendRequest(c *Client, url *url.URL, fn func(r *Request) (*Response, error)) (result *Result) {
	req, err := c.NewRequest(url.String())
	if err != nil {
		result = newResult(nil, err)
		return
	}

	resp, err := fn(req)
	result = newResult(resp, err)

	return
}

func (c *Client) addAuthenticationHeaders(header http.Header) {
	if c.AuthMethod != nil {
		header.Add("Authorization", c.AuthMethod.String())
	}

	if basicAuth, ok := c.AuthMethod.(BasicAuth); ok && basicAuth.OneTimePassword != "" {
		header.Add("X-GitHub-OTP", basicAuth.OneTimePassword)
	}
}
