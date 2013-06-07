package octokit

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestGet(t *testing.T) {
	c := NewClient()
	body, err := c.get("repos/jingweno/gh/commits", nil)

	assert.Equal(t, nil, err)
	assert.T(t, len(body) > 0)
}
