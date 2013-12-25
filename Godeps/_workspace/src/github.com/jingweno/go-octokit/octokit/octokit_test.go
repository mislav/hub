package octokit

import (
	"fmt"
	"github.com/bmizerany/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"testing"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the GitHub client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// setup sets up a test HTTP server along with a octokit.Client that is
// configured to talk to that test server.  Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// octokit client configured to use test server
	client = NewClientWith(server.URL, nil, TokenAuth{AccessToken: "token"})
}

// teardown closes the test HTTP server.
func tearDown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	assert.Equal(t, want, r.Method)
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	assert.Equal(t, want, r.Header.Get(header))
}

func testBody(t *testing.T, r *http.Request, want string) {
	body, _ := ioutil.ReadAll(r.Body)
	assert.Equal(t, want, string(body))
}

func respondWithJSON(w http.ResponseWriter, s string) {
	header := w.Header()
	header.Set("Content-Type", "application/json")
	respondWith(w, s)
}

func respondWithStatus(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func respondWith(w http.ResponseWriter, s string) {
	fmt.Fprint(w, s)
}

func testURLOf(path string) *url.URL {
	u, _ := url.ParseRequestURI(testURLStringOf(path))
	return u
}

func testURLStringOf(path string) string {
	return fmt.Sprintf("%s/%s", server.URL, path)
}

func loadFixture(f string) string {
	pwd, _ := os.Getwd()
	p := path.Join(pwd, "..", "fixtures", f)
	c, _ := ioutil.ReadFile(p)
	return string(c)
}
