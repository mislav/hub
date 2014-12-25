package octokit

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseError_empty_body(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		respondWith(w, "")
	})

	req, _ := client.NewRequest("error")
	_, err := req.Get(nil)
	assert.Contains(t, err.Error(), "400 - Problems parsing error message: EOF")

	e := err.(*ResponseError)
	assert.Equal(t, ErrorBadRequest, e.Type)
}

func TestResponseError_Error_400(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		respondWith(w, `{"message":"Problems parsing JSON"}`)
	})

	req, _ := client.NewRequest("error")
	_, err := req.Get(nil)
	assert.Contains(t, err.Error(), "400 - Problems parsing JSON")

	e := err.(*ResponseError)
	assert.Equal(t, ErrorBadRequest, e.Type)
}

func TestResponseError_Error_401(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		respondWith(w, `{"message":"Unauthorized"}`)
	})

	req, _ := client.NewRequest("error")
	_, err := req.Get(nil)
	assert.Contains(t, err.Error(), "401 - Unauthorized")

	e := err.(*ResponseError)
	assert.Equal(t, ErrorUnauthorized, e.Type)

	mux.HandleFunc("/error_2fa", func(w http.ResponseWriter, r *http.Request) {
		head := w.Header()
		head.Set("Content-Type", "application/json")
		head.Set("X-GitHub-OTP", "required; app")
		w.WriteHeader(http.StatusUnauthorized)
		respondWith(w, `{"message":"Unauthorized"}`)
	})

	req, _ = client.NewRequest("error_2fa")
	_, err = req.Get(nil)
	assert.Contains(t, err.Error(), "401 - Unauthorized")

	e = err.(*ResponseError)
	assert.Equal(t, ErrorOneTimePasswordRequired, e.Type)
}

func TestResponseError_Error_422_error(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(422)
		respondWith(w, `{"error":"No repository found for hubtopic"}`)
	})

	req, _ := client.NewRequest("error")
	_, err := req.Get(nil)
	assert.Contains(t, err.Error(), "Error: No repository found for hubtopic")

	e := err.(*ResponseError)
	assert.Equal(t, ErrorUnprocessableEntity, e.Type)
}

func TestResponseError_Error_422_error_summary(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(422)
		respondWith(w, `{"message":"Validation Failed", "errors": [{"resource":"Issue", "field": "title", "code": "missing_field"}]}`)
	})

	req, _ := client.NewRequest("error")
	_, err := req.Get(nil)
	assert.Contains(t, err.Error(), "422 - Validation Failed")
	assert.Contains(t, err.Error(), "missing_field error caused by title field on Issue resource")

	e := err.(*ResponseError)
	assert.Equal(t, ErrorUnprocessableEntity, e.Type)
}

func TestResponseError_Error_415(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		respondWith(w, `{"message":"Unsupported Media Type", "documentation_url":"http://developer.github.com/v3"}`)
	})

	req, _ := client.NewRequest("error")
	_, err := req.Get(nil)
	assert.Contains(t, err.Error(), "415 - Unsupported Media Type")
	assert.Contains(t, err.Error(), "// See: http://developer.github.com/v3")

	e := err.(*ResponseError)
	assert.Equal(t, ErrorUnsupportedMediaType, e.Type)
}

func TestResponseError_403_BadCredentials(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		respondWith(w, `{"message":"Bad credentials", "documentation_url":"https://developer.github.com/v3/"}`)
	})

	req, _ := client.NewRequest("error")
	_, err := req.Get(nil)
	assert.Contains(t, err.Error(), "403 - Bad credentials")
	assert.Contains(t, err.Error(), "// See: https://developer.github.com/v3/")

	e := err.(*ResponseError)
	assert.Equal(t, ErrorForbidden, e.Type)
}

func TestResponseError_403_RateLimit(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		respondWith(w, `{"message":"API rate limit exceeded for 1.2.3.4", "documentation_url":"https://developer.github.com/v3/#rate-limiting"}`)
	})

	req, _ := client.NewRequest("error")
	_, err := req.Get(nil)
	assert.Contains(t, err.Error(), "403 - API rate limit exceeded for 1.2.3.4")
	assert.Contains(t, err.Error(), "// See: https://developer.github.com/v3/#rate-limiting")

	e := err.(*ResponseError)
	assert.Equal(t, ErrorTooManyRequests, e.Type)
}

func TestResponseError_403_LoginLimit(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		respondWith(w, `{"message":"login attempts exceeded", "documentation_url":"https://developer.github.com/v3/"}`)
	})

	req, _ := client.NewRequest("error")
	_, err := req.Get(nil)
	assert.Contains(t, err.Error(), "403 - login attempts exceeded")
	assert.Contains(t, err.Error(), "// See: https://developer.github.com/v3/")

	e := err.(*ResponseError)
	assert.Equal(t, ErrorTooManyLoginAttempts, e.Type)
}
