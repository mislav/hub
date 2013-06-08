package octokit

import (
	"github.com/bmizerany/assert"
	"os"
	"testing"
)

func TestAuthorizations(t *testing.T) {
	c := NewClientWithPassword(os.Getenv("GITHUB_USER"), os.Getenv("GITHUB_PASSWORD"))
	auths, err := c.Authorizations()

	assert.Equal(t, nil, err)
	assert.T(t, len(auths) > 0)
}
