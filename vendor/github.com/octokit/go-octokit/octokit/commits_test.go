package octokit

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommitsService_all(t *testing.T) {
	setup()
	defer tearDown()

	//Username URL
	stubGet(t, "/repos/octocat/Hello-World/commits", "commits", nil)

	commits, result := client.Commits().All(&CommitsURL, M{
		"owner": "octocat",
		"repo":  "Hello-World",
	})

	assert.False(t, result.HasError())
	assert.Equal(t, "6dcb09b5b57875f334f61aebed695e2e4193db5e", commits[0].Sha)
	assert.Equal(t, "https://github.com/images/error/octocat_happy.gif", commits[0].Author.AvatarURL)
	assert.Len(t, commits, 1)

	commit := commits[0].Commit
	assert.Equal(t, "Monalisa Octocat", commit.Author.Name)
	assert.Equal(t, "support@github.com", commit.Committer.Email)
	assert.Equal(t, "Fix all the bugs", commit.Message)

	//Nil case
	commitsNil, resultNil := client.Commits().All(nil, M{
		"owner": "octocat",
		"repo":  "Hello-World",
	})

	assert.False(t, resultNil.HasError())
	assert.Equal(t, commitsNil, commits)

	//Error case
	var invalid = Hyperlink("{")
	commitsErr, resultErr := client.Commits().All(&invalid, M{})
	assert.True(t, resultErr.HasError())
	assert.Equal(t, commitsErr, make([]Commit, 0))
}

func TestCommitsService_One(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/repos/octokit/go-octokit/commits/4351fb69b8d5ed075e9cd844e67ad2114b335c82", "commit", nil)

	commit, result := client.Commits().One(&CommitsURL, M{
		"owner": "octokit",
		"repo":  "go-octokit",
		"sha":   "4351fb69b8d5ed075e9cd844e67ad2114b335c82",
	})

	assert.False(t, result.HasError())
	assert.Equal(t, "4351fb69b8d5ed075e9cd844e67ad2114b335c82", commit.Sha)
	assert.Equal(t, "https://api.github.com/repos/octokit/go-octokit/commits/4351fb69b8d5ed075e9cd844e67ad2114b335c82", commit.URL)

	files := commit.Files
	assert.Len(t, files, 35)

	commitNil, resultNil := client.Commits().One(nil, M{
		"owner": "octokit",
		"repo":  "go-octokit",
		"sha":   "4351fb69b8d5ed075e9cd844e67ad2114b335c82",
	})
	assert.False(t, resultNil.HasError())
	assert.Equal(t, commit, commitNil)

	//Error case
	var invalid = Hyperlink("{")
	commitErr, resultErr := client.Commits().One(&invalid, M{})
	assert.True(t, resultErr.HasError())
	assert.Equal(t, commitErr, Commit{})
}

func TestCommitsService_Patch(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/octokit/go-octokit/commits/b6d21008bf7553a29ad77ee0a8bb3b66e6f11aa2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", patchMediaType)
		respondWithJSON(w, loadFixture("commit.patch"))
	})

	patch, result := client.Commits().Patch(&CommitsURL, M{
		"owner": "octokit",
		"repo":  "go-octokit",
		"sha":   "b6d21008bf7553a29ad77ee0a8bb3b66e6f11aa2",
	})

	assert.False(t, result.HasError())
	content, err := ioutil.ReadAll(patch)
	assert.NoError(t, err)
	assert.NotEmpty(t, content)

	patchNil, resultNil := client.Commits().Patch(nil, M{
		"owner": "octokit",
		"repo":  "go-octokit",
		"sha":   "b6d21008bf7553a29ad77ee0a8bb3b66e6f11aa2",
	})
	assert.False(t, resultNil.HasError())
	contentNil, errNil := ioutil.ReadAll(patchNil)
	assert.NoError(t, errNil)
	assert.Equal(t, content, contentNil)

	//Error case
	var invalid = Hyperlink("{")
	commitErr, resultErr := client.Commits().Patch(&invalid, M{})
	assert.True(t, resultErr.HasError())
	assert.Nil(t, commitErr)
}
