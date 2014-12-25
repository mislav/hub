package octokit

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootService_One(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("root.json"))
	})

	url, err := RootURL.Expand(nil)
	assert.NoError(t, err)

	root, result := client.Root(url).One()
	assert.False(t, result.HasError())
	assert.Equal(t, "https://api.github.com/users/{user}", string(root.UserURL))
}

func TestClientRel(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("root.json"))
	})

	u, err := client.Rel("user", M{"user": "root"})
	assert.NoError(t, err)
	assert.Equal(t, "https://api.github.com/users/root", u.String())
}
