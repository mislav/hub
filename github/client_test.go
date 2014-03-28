package github

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/octokit/go-octokit/octokit"
)

func TestClient_FormatError(t *testing.T) {
	e := &octokit.ResponseError{
		Response: &http.Response{
			StatusCode: 401,
			Status:     "401 Not Found",
		},
	}

	err := FormatError("action", e)
	assert.Equal(t, "Error action: Not Found (HTTP 401)", fmt.Sprintf("%s", err))

	e = &octokit.ResponseError{
		Response: &http.Response{
			StatusCode: 422,
			Status:     "422 Unprocessable Entity",
		},
		Message: "error message",
	}
	err = FormatError("action", e)
	assert.Equal(t, "Error action: Unprocessable Entity (HTTP 422)\nerror message", fmt.Sprintf("%s", err))
}

func TestClient_warnExistenceOfRepo(t *testing.T) {
	project := &Project{
		Name:  "hub",
		Owner: "github",
		Host:  "github.com",
	}
	e := &octokit.ResponseError{
		Response: &http.Response{
			StatusCode: 404,
			Status:     "404 Not Found",
		},
		Message: "error message",
	}

	err := warnExistenceOfRepo(project, e)
	assert.Equal(t, "Are you sure that github.com/github/hub exists?", fmt.Sprintf("%s", err))
}
