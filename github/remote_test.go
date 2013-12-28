package github

import (
	"github.com/bmizerany/assert"
	"strings"
	"testing"
)

func TestOriginRemote(t *testing.T) {
	gitRemote, _ := OriginRemote()
	assert.Equal(t, "origin", gitRemote.Name)
	assert.T(t, strings.Contains(gitRemote.URL.String(), "gh"))
}
