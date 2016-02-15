package octokit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchService_Users(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/search/users", "user_search", nil)

	searchResults, result := client.Search().Users(nil, M{"query": "dhruvsinghal"})

	assert.False(t, result.HasError())
	assert.False(t, searchResults.IncompleteResults)
	assert.Equal(t, 2, searchResults.TotalCount)
	assert.Equal(t, 2, len(searchResults.Items))
	assert.Equal(t, 3338555, searchResults.Items[0].ID)
	assert.EqualValues(t, "dhruvsinghal", searchResults.Items[0].Login)
	assert.Equal(t, 9294878, searchResults.Items[1].ID)
	assert.EqualValues(t, "dhruvsinghal5", searchResults.Items[1].Login)
}

func TestSearchService_Issues(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/search/issues", "issue_search", nil)

	searchResults, result := client.Search().Issues(nil, M{"query": "color"})

	assert.False(t, result.HasError())
	assert.False(t, searchResults.IncompleteResults)
	assert.Equal(t, 180172, searchResults.TotalCount)
	assert.Equal(t, 1551, searchResults.Items[0].Number)
	assert.EqualValues(t, "Colorizer", searchResults.Items[0].Title)
	assert.Equal(t, 3402, searchResults.Items[1].Number)
	assert.EqualValues(t, "Colorizer", searchResults.Items[1].Title)
}

func TestSearchService_Repositories(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/search/repositories", "repository_search", nil)

	searchResults, result := client.Search().Repositories(nil,
		M{"query": "asdfghjk"})

	assert.False(t, result.HasError())
	assert.False(t, searchResults.IncompleteResults)
	assert.Equal(t, 21, searchResults.TotalCount)
	assert.Equal(t, 21, len(searchResults.Items))
	assert.Equal(t, 8269299, searchResults.Items[0].ID)
	assert.EqualValues(t, "ysai/asdfghjk", searchResults.Items[0].FullName)
	assert.Equal(t, 8511889, searchResults.Items[1].ID)
	assert.EqualValues(t, "ines949494/ikadasd", searchResults.Items[1].FullName)
}

func TestSearchService_Code(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/search/code", "code_search", nil)

	searchResults, result := client.Search().Code(nil, M{
		"query": "addClass in:file language:js repo:jquery/jquery"})

	assert.False(t, result.HasError())
	assert.False(t, searchResults.IncompleteResults)
	assert.Equal(t, 7, searchResults.TotalCount)
	assert.Equal(t, 7, len(searchResults.Items))
	assert.EqualValues(t, "classes.js", searchResults.Items[0].Name)
	assert.EqualValues(t, "src/attributes/classes.js", searchResults.Items[0].Path)
	assert.EqualValues(t,
		"f9dba94f7de43d6b6b7256e05e0d17c4741a4cde", searchResults.Items[0].SHA)
	assert.EqualValues(t,
		"https://api.github.com/repositories/167174/contents/src/attributes/classes.js?ref=53aa87f3bf4284763405f3eb8affff296e55ba4f", string(searchResults.Items[0].URL))
	assert.EqualValues(t,
		"https://api.github.com/repositories/167174/git/blobs/f9dba94f7de43d6b6b7256e05e0d17c4741a4cde", searchResults.Items[0].GitURL)
	assert.EqualValues(t,
		"https://github.com/jquery/jquery/blob/53aa87f3bf4284763405f3eb8affff296e55ba4f/src/attributes/classes.js", searchResults.Items[0].HTMLURL)
	assert.Equal(t, 167174, searchResults.Items[0].Repository.ID)
	assert.EqualValues(t,
		"jquery/jquery", searchResults.Items[0].Repository.FullName)
}
