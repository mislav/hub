package octokit

import (
	"github.com/bmizerany/assert"
	"os"
	"testing"
)

func TestStatuses(t *testing.T) {
	c := NewClientWithToken(os.Getenv("GITHUB_TOKEN"))
	repo := Repository{"gh", "jingweno"}

	statuses, err := c.Statuses(repo, "99b0f36b24e25a644ed70ace601da857eea4cf72")
	assert.Equal(t, nil, err)
	assert.T(t, len(statuses) == 0)
}
