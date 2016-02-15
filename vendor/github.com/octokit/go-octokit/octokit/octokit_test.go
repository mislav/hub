package octokit

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the GitHub client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// A http.Transport subtype that re-routes all requests in testing to the local
// server as indicated by `overrideURL`.
type TestTransport struct {
	http.RoundTripper
	overrideURL *url.URL
}

func (t TestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = cloneRequest(req)
	req.Header.Set("X-Original-Scheme", req.URL.Scheme)
	req.URL.Scheme = t.overrideURL.Scheme
	req.URL.Host = t.overrideURL.Host
	return t.RoundTripper.RoundTrip(req)
}

func cloneRequest(r *http.Request) *http.Request {
	r2 := new(http.Request)
	*r2 = *r
	r2.URL, _ = url.Parse(r.URL.String())
	r2.Header = make(http.Header)
	for k, s := range r.Header {
		r2.Header[k] = s
	}
	return r2
}

// setup sets up a test HTTP server along with a octokit.Client that is
// configured to talk to that test server.  Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	serverURL, _ := url.Parse(server.URL)

	httpClient := http.Client{
		Transport: TestTransport{
			RoundTripper: http.DefaultTransport,
			overrideURL:  serverURL,
		},
	}

	// octokit client configured to use test server
	client = NewClientWith(
		gitHubAPIURL,
		userAgent,
		TokenAuth{AccessToken: "token"},
		&httpClient,
	)
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
	assert.EqualValues(t, want, string(body))
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
	p := filepath.Join(pwd, "..", "fixtures", f)
	c, _ := ioutil.ReadFile(p)
	return string(c)
}

func stubGet(t *testing.T, path string, fixture string,
	respHeaderParams map[string]string) {
	httpTestHelper(t, path, fixture, "GET", nil, "", respHeaderParams, 0)
}

func stubGetwCode(t *testing.T, path string, fixture string,
	respHeaderParams map[string]string, respStatusCode int) {
	httpTestHelper(t, path, fixture, "GET",
		nil, "", respHeaderParams, respStatusCode)
}

func stubPost(t *testing.T, path string, fixture string,
	wantReqHeader map[string]string, wantReqBody string,
	respHeaderParams map[string]string) {
	httpTestHelper(t, path, fixture, "POST",
		wantReqHeader, wantReqBody, respHeaderParams, 0)
}

func stubPostwCode(t *testing.T, path string, fixture string,
	wantReqHeader map[string]string, wantReqBody string,
	respHeaderParams map[string]string, respStatusCode int) {
	httpTestHelper(t, path, fixture, "POST",
		wantReqHeader, wantReqBody, respHeaderParams, respStatusCode)
}

func stubPutwCode(t *testing.T, path string, fixture string,
	wantReqHeader map[string]string, wantReqBody string,
	respHeaderParams map[string]string, respStatusCode int) {
	httpTestHelper(t, path, fixture, "PUT",
		wantReqHeader, wantReqBody, respHeaderParams, respStatusCode)
}

func stubPatch(t *testing.T, path string, fixture string,
	wantReqHeader map[string]string, wantReqBody string,
	respHeaderParams map[string]string) {
	httpTestHelper(t, path, fixture, "PATCH",
		wantReqHeader, wantReqBody, respHeaderParams, 0)
}

func stubDeletewCode(t *testing.T, path string,
	respHeaderParams map[string]string, respStatusCode int) {
	httpTestHelper(t, path, "", "DELETE",
		nil, "", respHeaderParams, respStatusCode)
}

func stubDeletewCodewBody(t *testing.T, path string, wantReqBody string,
	respHeaderParams map[string]string, respStatusCode int) {
	httpTestHelper(t, path, "", "DELETE",
		nil, wantReqBody, respHeaderParams, respStatusCode)
}

func httpTestHelper(t *testing.T, path string, fixture string,
	wantReqMethod string, wantReqHeader map[string]string, wantReqBody string,
	respHeaderParams map[string]string, respStatusCode int) {

	if mux == nil {
		panic(fmt.Errorf("test HTTP server has not been set up"))
	}
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		// Testing the incoming request
		testMethod(t, r, wantReqMethod)
		if wantReqHeader != nil {
			for k, v := range wantReqHeader {
				testHeader(t, r, k, v)
			}
		}

		if wantReqBody != "" {
			testBody(t, r, wantReqBody)
		}

		// Construct the outgoing response
		if respHeaderParams != nil {
			header := w.Header()
			for k, v := range respHeaderParams {
				header.Set(k, v)
			}
		}
		if respStatusCode > 0 {
			w.WriteHeader(respStatusCode)
		}
		if fixture != "" {
			respondWithJSON(w, loadFixture(fixture+".json"))
		}
	})
}
