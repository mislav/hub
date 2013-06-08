package octokit

import (
	"github.com/bmizerany/assert"
	"os"
	"testing"
)

func TestAuthenticatedUser(t *testing.T) {
	c := NewClientWithToken(os.Getenv("GITHUB_TOKEN"))
	user, err := c.AuthenticatedUser()

	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", user.Login)
}
