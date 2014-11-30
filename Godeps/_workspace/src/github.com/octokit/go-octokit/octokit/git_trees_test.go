package octokit

import (
	"net/http"
	"testing"

	"github.com/bmizerany/assert"
)

func TestGitTreesService_One(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/pengwynn/flint/git/trees/master", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("tree.json"))
	})

	url, err := GitTreesURL.Expand(M{
		"owner": "pengwynn",
		"repo":  "flint",
		"sha":   "master",
	})
	assert.Equal(t, nil, err)
	tree, result := client.GitTrees(url).One()

	assert.T(t, !result.HasError())
	assert.Equal(t, "9c1337e761bbd517f3cc1b5acb9373b17f4810e8", tree.Sha)
	assert.Equal(t, "https://api.github.com/repos/pengwynn/flint/git/trees/9c1337e761bbd517f3cc1b5acb9373b17f4810e8", tree.URL)

	entries := tree.Tree
	assert.Equal(t, 9, len(entries))
}

func TestGitTreesService_One_Recursive(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/pengwynn/flint/git/trees/master", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("tree_recursive.json"))
	})

	url, err := GitTreesURL.Expand(M{
		"owner":     "pengwynn",
		"repo":      "flint",
		"sha":       "master",
		"recursive": "1",
	})
	assert.Equal(t, nil, err)
	tree, result := client.GitTrees(url).One()

	assert.T(t, !result.HasError())
	assert.Equal(t, "9c1337e761bbd517f3cc1b5acb9373b17f4810e8", tree.Sha)
	assert.Equal(t, "https://api.github.com/repos/pengwynn/flint/git/trees/9c1337e761bbd517f3cc1b5acb9373b17f4810e8", tree.URL)

	entries := tree.Tree
	assert.Equal(t, 15, len(entries))
}
