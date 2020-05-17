package github

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/github/hub/v2/internal/assert"
)

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

	assert.Equal(t, "hub for <unidentified machine>", note)

	note, err = authTokenNote(2)
	assert.Equal(t, nil, err)
	assert.Equal(t, "hub for <unidentified machine> 2", note)

	os.Setenv("HUB_MACHINE", "mydevmachine")

	note, err = authTokenNote(1)
	assert.Equal(t, nil, err)
	assert.Equal(t, "hub for mydevmachine", note)

	note, err = authTokenNote(2)
	assert.Equal(t, nil, err)
	assert.Equal(t, "hub for mydevmachine 2", note)
}
