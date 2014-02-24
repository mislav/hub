package github

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/octokit/go-octokit/octokit"
)

func TestClient_ApiEndpoint(t *testing.T) {
	gh := &Client{Credentials: &Credentials{Host: "github.com"}}
	assert.Equal(t, "https://api.github.com", gh.apiEndpoint())

	gh = &Client{Credentials: &Credentials{Host: "github.corporate.com"}}
	assert.Equal(t, "https://github.corporate.com", gh.apiEndpoint())

	gh = &Client{Credentials: &Credentials{Host: "http://github.corporate.com"}}
	assert.Equal(t, "http://github.corporate.com", gh.apiEndpoint())
}

func TestClient_formatError(t *testing.T) {
	result := &octokit.Result{
		Err: &octokit.ResponseError{
			Response: &http.Response{
				StatusCode: 401,
				Status:     "401 Not Found",
			},
		},
	}
	err := formatError("action", result)
	assert.Equal(t, "Error action: Not Found (HTTP 401)", fmt.Sprintf("%s", err))

	result = &octokit.Result{
		Err: &octokit.ResponseError{
			Response: &http.Response{
				StatusCode: 422,
				Status:     "422 Unprocessable Entity",
			},
			Message: "error message",
		},
	}
	err = formatError("action", result)
	assert.Equal(t, "Error action: Unprocessable Entity (HTTP 422)\nerror message", fmt.Sprintf("%s", err))
}

func TestClient_warnExistenceOfRepo(t *testing.T) {
	project := &Project{
		Name:  "hub",
		Owner: "github",
		Host:  "github.com",
	}
	result := &octokit.Result{
		Err: &octokit.ResponseError{
			Response: &http.Response{
				StatusCode: 404,
				Status:     "404 Not Found",
			},
			Message: "error message",
		},
	}

	err := warnExistenceOfRepo(project, result)
	assert.Equal(t, "Are you sure that github.com/github/hub exists?", fmt.Sprintf("%s", err))
}
