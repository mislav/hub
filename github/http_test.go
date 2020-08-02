package github

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/github/hub/v2/internal/assert"
)

func setupTestServer(unixSocket string) *testServer {
	m := http.NewServeMux()
	s := httptest.NewServer(m)
	u, _ := url.Parse(s.URL)

	if unixSocket != "" {
		os.Remove(unixSocket)
		unixListener, err := net.Listen("unix", unixSocket)
		if err != nil {
			log.Fatal("Unable to listen on unix-socket: ", err)
		}
		go http.Serve(unixListener, m)
	}

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
	s := setupTestServer("")
	defer s.Close()

	s.HandleFunc("/override", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "https", r.Header.Get("X-Original-Scheme"))
		assert.Equal(t, "example.com", r.Host)
	})

	c := newHTTPClient(s.URL.String(), false, "")
	c.Get("https://example.com/override")

	s.HandleFunc("/not-override", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "", r.Header.Get("X-Original-Scheme"))
		assert.Equal(t, s.URL.Host, r.Host)
	})

	c = newHTTPClient("", false, "")
	c.Get(fmt.Sprintf("%s/not-override", s.URL.String()))
}

func TestNewHttpClient_UnixSocket(t *testing.T) {
	sock := "/tmp/hub-go.sock"
	s := setupTestServer(sock)
	defer s.Close()

	s.HandleFunc("/unix-socket", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("unix-socket-works"))
	})
	c := newHTTPClient("", false, sock)
	resp, err := c.Get(fmt.Sprintf("%s/unix-socket", s.URL.String()))
	assert.Equal(t, nil, err)
	result, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "unix-socket-works", string(result))
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
