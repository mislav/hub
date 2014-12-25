package octokit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHyperlink_Expand(t *testing.T) {
	link := Hyperlink("https://api.github.com/users/{user}")
	url, err := link.Expand(M{"user": "jingweno"})
	assert.NoError(t, err)
	assert.Equal(t, "https://api.github.com/users/jingweno", url.String())

	link = Hyperlink("https://api.github.com/user")
	url, err = link.Expand(nil)
	assert.NoError(t, err)
	assert.Equal(t, "https://api.github.com/user", url.String())

	url, err = link.Expand(M{})
	assert.NoError(t, err)
	assert.Equal(t, "https://api.github.com/user", url.String())
}
