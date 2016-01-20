package octokit

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestHyperlink_Expand(t *testing.T) {
	link := Hyperlink("https://api.github.com/users/{user}")
	url, err := link.Expand(M{"user": "jingweno"})
	assert.Equal(t, nil, err)
	assert.Equal(t, "https://api.github.com/users/jingweno", url.String())

	link = Hyperlink("https://api.github.com/user")
	url, err = link.Expand(nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, "https://api.github.com/user", url.String())

	url, err = link.Expand(M{})
	assert.Equal(t, nil, err)
	assert.Equal(t, "https://api.github.com/user", url.String())
}
