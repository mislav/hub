package github

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bmizerany/assert"
)

func setupTestServer() *testServer {
	m := http.NewServeMux()
	s := httptest.NewServer(m)
	u, _ := url.Parse(s.URL)

	return &testServer{
		Server:   s,
		ServeMux: m,
		URL:      u,
	}
}

type testServer struct {
	*http.ServeMux
	Server *httptest.Server
	URL    *url.URL
}

func (s *testServer) Close() {
	s.Server.Close()
}

func TestNewHttpClient_OverrideURL(t *testing.T) {
	s := setupTestServer()
	defer s.Close()

	s.HandleFunc("/override", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "https", r.Header.Get("X-Original-Scheme"))
		assert.Equal(t, "example.com", r.Host)
	})

	c := newHttpClient(s.URL.String(), false)
	c.Get("https://example.com/override")

	s.HandleFunc("/not-override", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "", r.Header.Get("X-Original-Scheme"))
		assert.Equal(t, s.URL.Host, r.Host)
	})

	c = newHttpClient("", false)
	c.Get(fmt.Sprintf("%s/not-override", s.URL.String()))
}

func TestVerboseTransport_VerbosePrintln(t *testing.T) {
	var b bytes.Buffer
	tr := &verboseTransport{
		Out:       &b,
		Colorized: true,
	}

	tr.verbosePrintln("foo")
	assert.Equal(t, "\033[36mfoo\033[0m\n", b.String())
}
