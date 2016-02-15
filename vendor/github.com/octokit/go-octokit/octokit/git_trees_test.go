package octokit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitTreesService_One(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/repos/pengwynn/flint/git/trees/master", "tree", nil)

	url, err := GitTreesURL.Expand(M{
		"owner": "pengwynn",
		"repo":  "flint",
		"sha":   "master",
	})
	assert.NoError(t, err)
	tree, result := client.GitTrees(url).One()

	assert.False(t, result.HasError())
	assert.Equal(t, "9c1337e761bbd517f3cc1b5acb9373b17f4810e8", tree.Sha)
	assert.Equal(t, "https://api.github.com/repos/pengwynn/flint/git/trees/9c1337e761bbd517f3cc1b5acb9373b17f4810e8", tree.URL)

	entries := tree.Tree
	assert.Len(t, entries, 9)
}

func TestGitTreesService_One_Recursive(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/repos/pengwynn/flint/git/trees/master", "tree_recursive", nil)

	url, err := GitTreesURL.Expand(M{
		"owner":     "pengwynn",
		"repo":      "flint",
		"sha":       "master",
		"recursive": "1",
	})
	assert.NoError(t, err)
	tree, result := client.GitTrees(url).One()

	assert.False(t, result.HasError())
	assert.Equal(t, "9c1337e761bbd517f3cc1b5acb9373b17f4810e8", tree.Sha)
	assert.Equal(t, "https://api.github.com/repos/pengwynn/flint/git/trees/9c1337e761bbd517f3cc1b5acb9373b17f4810e8", tree.URL)

	entries := tree.Tree
	assert.Len(t, entries, 15)
}
