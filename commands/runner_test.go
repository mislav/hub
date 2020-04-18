package commands

import (
	"testing"

	"github.com/github/hub/v2/internal/assert"
)

func TestRunner_splitAliasCmd(t *testing.T) {
	_, err := splitAliasCmd("!source ~/.zshrc")
	assert.NotEqual(t, nil, err)

	words, err := splitAliasCmd("log --pretty=oneline --abbrev-commit --graph --decorate")
	assert.Equal(t, nil, err)
	assert.Equal(t, 5, len(words))

	_, err = splitAliasCmd("")
	assert.NotEqual(t, nil, err)
}
