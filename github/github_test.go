package github

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestApiEndpoint(t *testing.T) {
	gh := &GitHub{Project: &Project{Name: "foo", Owner: "bar", Host: "github.com"}}
	assert.Equal(t, "https://api.github.com", gh.apiEndpoint())

	gh = &GitHub{Project: &Project{Name: "foo", Owner: "bar", Host: "github.corporate.com"}}
	assert.Equal(t, "https://github.corporate.com", gh.apiEndpoint())

	gh = &GitHub{Project: &Project{Name: "foo", Owner: "bar", Host: "http://github.corporate.com"}}
	assert.Equal(t, "http://github.corporate.com", gh.apiEndpoint())
}
