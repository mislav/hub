package github

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/bmizerany/assert"
)

var Hub1Re = regexp.MustCompile("hub for (.+)@(.+)")

var Hub2Re = regexp.MustCompile("hub for (.+)@(.+) 2")

func TestClient_FormatError(t *testing.T) {
	e := &errorInfo{
		Response: &http.Response{
			StatusCode: 401,
			Status:     "401 Not Found",
		},
	}

	err := FormatError("action", e)
	assert.Equal(t, "Error action: Not Found (HTTP 401)", fmt.Sprintf("%s", err))

	e = &errorInfo{
		Response: &http.Response{
			StatusCode: 422,
			Status:     "422 Unprocessable Entity",
		},
		Message: "error message",
	}
	err = FormatError("action", e)
	assert.Equal(t, "Error action: Unprocessable Entity (HTTP 422)\nerror message", fmt.Sprintf("%s", err))
}

func TestAuthTokenNote(t *testing.T) {
	note, err := authTokenNote(1)
	assert.Equal(t, nil, err)

	assert.T(t, Hub1Re.MatchString(note))

	note, err = authTokenNote(2)
	assert.Equal(t, nil, err)

	assert.T(t, Hub2Re.MatchString(note))

}
