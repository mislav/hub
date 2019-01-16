package commands

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestRunner_splitAliasCmd(t *testing.T) {
	words, err := splitAliasCmd("!source ~/.zshrc")
	assert.NotEqual(t, nil, err)

	words, err = splitAliasCmd("log --pretty=oneline --abbrev-commit --graph --decorate")
	assert.Equal(t, nil, err)
	assert.Equal(t, 5, len(words))

	words, err = splitAliasCmd("")
	assert.NotEqual(t, nil, err)
}
