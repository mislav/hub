package github

import (
	"fmt"
	"github.com/github/hub/git"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/bmizerany/assert"
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

	reg := regexp.MustCompile("hub for (.+)@(.+)")
	assert.T(t, reg.MatchString(note))

	note, err = authTokenNote(2)
	assert.Equal(t, nil, err)

	reg = regexp.MustCompile("hub for (.+)@(.+) 2")
	assert.T(t, reg.MatchString(note))

}

func TestUserAgent(t *testing.T){
	t.Run("set from env", func(t *testing.T){
		os.Setenv("HUB_USERAGENT","git/curl")
		SetUserAgent()
		assert.Equal(t, "git/curl", UserAgent)
	})

	t.Run("set from git", func(t *testing.T) {
		dir, _ := ioutil.TempDir("","foo.bar")
		defer os.RemoveAll(dir)
		os.Chdir(dir)

		git.Run("init")
		git.Run("config","hub.useragent","foo.bar")

		SetUserAgent()
		assert.Equal(t, "foo.bar", UserAgent)
	})

}
