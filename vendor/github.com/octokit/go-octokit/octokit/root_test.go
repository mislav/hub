package octokit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootService_One(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/", "root", nil)

	url, err := RootURL.Expand(nil)
	assert.NoError(t, err)

	root, result := client.Root(url).One()
	assert.False(t, result.HasError())
	assert.Equal(t, "https://api.github.com/users/{user}", string(root.UserURL))
}

func TestClientRel(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/", "root", nil)

	u, err := client.Rel("user", M{"user": "root"})
	assert.NoError(t, err)
	assert.Equal(t, "https://api.github.com/users/root", u.String())
}
