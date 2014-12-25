package sawyer

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/jingweno/go-sawyer/mediatype"
)

// The default httpClient used if one isn't specified.
var httpClient = &http.Client{}

// A Client wraps an *http.Client with a base url Endpoint and common header and
// query values.
type Client struct {
	HttpClient *http.Client
	Endpoint   *url.URL
	Header     http.Header
	Query      url.Values
}

// New returns a new Client with a given a URL and an optional client.
func New(endpoint *url.URL, client *http.Client) *Client {
	if client == nil {
		client = httpClient
	}

	if len(endpoint.Path) > 0 && !strings.HasSuffix(endpoint.Path, "/") {
		endpoint.Path = endpoint.Path + "/"
	}

	return &Client{client, endpoint, make(http.Header), endpoint.Query()}
}

// NewFromString returns a new Client given a string URL and an optional client.
func NewFromString(endpoint string, client *http.Client) (*Client, error) {
	e, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	return New(e, client), nil
}

// ResolveReference resolves a URI reference to an absolute URI from an absolute
// base URI.  It also merges the query values.
func (c *Client) ResolveReference(u *url.URL) *url.URL {
	absurl := c.Endpoint.ResolveReference(u)
	if len(c.Query) > 0 {
		absurl.RawQuery = mergeQueries(c.Query, absurl.Query())
	}
	return absurl
}

// ResolveReference resolves a string URI reference to an absolute URI from an
// absolute base URI.  It also merges the query values.
func (c *Client) ResolveReferenceString(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	return c.ResolveReference(u).String(), nil
}

func mergeQueries(queries ...url.Values) string {
	merged := make(url.Values)
	for _, q := range queries {
		if len(q) == 0 {
			break
		}

		for key, _ := range q {
			merged.Set(key, q.Get(key))
		}
	}
	return merged.Encode()
}

func init() {
	mediatype.AddDecoder("json", func(r io.Reader) mediatype.Decoder {
		return json.NewDecoder(r)
	})
	mediatype.AddEncoder("json", func(w io.Writer) mediatype.Encoder {
		return json.NewEncoder(w)
	})
}
