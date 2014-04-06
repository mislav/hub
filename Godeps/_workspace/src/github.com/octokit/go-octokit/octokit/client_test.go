package octokit

import (
	"net/http"
	"testing"

	"github.com/bmizerany/assert"
)

func TestSuccessfulGet(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", defaultMediaType)
		testHeader(t, r, "User-Agent", userAgent)
		testHeader(t, r, "Authorization", "token token")
		respondWithJSON(w, `{"login": "octokit"}`)
	})

	req, err := client.NewRequest("foo")
	assert.Equal(t, nil, err)

	var output map[string]interface{}
	_, err = req.Get(&output)
	assert.Equal(t, nil, err)
	assert.Equal(t, "octokit", output["login"])
}

func TestSuccessfulGet_BasicAuth(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", defaultMediaType)
		testHeader(t, r, "User-Agent", userAgent)
		testHeader(t, r, "Authorization", "Basic amluZ3dlbm86cGFzc3dvcmQ=")
		testHeader(t, r, "X-GitHub-OTP", "OTP")
		respondWithJSON(w, `{"login": "octokit"}`)
	})

	client = NewClientWith(
		server.URL,
		userAgent,
		BasicAuth{
			Login:           "jingweno",
			Password:        "password",
			OneTimePassword: "OTP",
		},
		nil)
	req, err := client.NewRequest("foo")
	assert.Equal(t, nil, err)

	var output map[string]interface{}
	_, err = req.Get(&output)
	assert.Equal(t, nil, err)
	assert.Equal(t, "octokit", output["login"])
}

func TestGetWithoutDecoder(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		head := w.Header()
		head.Set("Content-Type", "application/booya+booya")
		respondWith(w, `{"login": "octokit"}`)
	})

	req, err := client.NewRequest("foo")
	assert.Equal(t, nil, err)

	var output map[string]interface{}
	_, err = req.Get(output)
	assert.NotEqual(t, nil, err)
}

func TestGetResponseError(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		respondWith(w, `{"message": "not found"}`)
	})

	req, err := client.NewRequest("foo")
	assert.Equal(t, nil, err)

	var output map[string]interface{}
	_, err = req.Get(output)
	assert.NotEqual(t, nil, err)
	respErr, ok := err.(*ResponseError)
	assert.Tf(t, ok, "should be able to convert to *ResponseError")
	assert.Equal(t, "not found", respErr.Message)
	assert.Equal(t, ErrorNotFound, respErr.Type)
}

func TestSuccessfulPost(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", defaultMediaType)
		testHeader(t, r, "Content-Type", defaultMediaType)
		testHeader(t, r, "User-Agent", userAgent)
		testHeader(t, r, "Authorization", "token token")
		testBody(t, r, "{\"input\":\"bar\"}\n")
		respondWithJSON(w, `{"login": "octokit"}`)
	})

	req, err := client.NewRequest("foo")
	assert.Equal(t, nil, err)

	input := map[string]interface{}{"input": "bar"}
	var output map[string]interface{}
	_, err = req.Post(input, &output)
	assert.Equal(t, nil, err)
	assert.Equal(t, "octokit", output["login"])
}

func TestAddHeader(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Foo", "Bar")
		assert.Equal(t, "example.com", r.Host)
		respondWithJSON(w, `{"login": "octokit"}`)
	})

	client.Header.Set("Host", "example.com")
	client.Header.Set("Foo", "Bar")
	req, err := client.NewRequest("foo")
	assert.Equal(t, nil, err)

	_, err = req.Get(nil)
	assert.Equal(t, nil, err)
}
