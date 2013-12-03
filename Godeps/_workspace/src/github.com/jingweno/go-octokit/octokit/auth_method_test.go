package octokit

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	basicAuth := BasicAuth{Login: "jingweno", Password: "password"}
	assert.Equal(t, "Basic amluZ3dlbm86cGFzc3dvcmQ=", basicAuth.String())
}

func TestTokenAuth(t *testing.T) {
	tokenAuth := TokenAuth{AccessToken: "token"}
	assert.Equal(t, "token token", tokenAuth.String())
}
