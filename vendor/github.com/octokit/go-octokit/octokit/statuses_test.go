package octokit

import (
	"net/http"
	"testing"

	"github.com/bmizerany/assert"
)

func TestStatuses(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/jingweno/gh/statuses/740211b9c6cd8e526a7124fe2b33115602fbc637", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("statuses.json"))
	})

	sha := "740211b9c6cd8e526a7124fe2b33115602fbc637"
	url, err := StatusesURL.Expand(M{"owner": "jingweno", "repo": "gh", "ref": sha})
	assert.Equal(t, nil, err)

	statuses, err := client.Statuses(url).All()

	assert.Equal(t, 2, len(statuses))
	firstStatus := statuses[0]
	assert.Equal(t, "pending", firstStatus.State)
	assert.Equal(t, "The Travis CI build is in progress", firstStatus.Description)
	assert.Equal(t, "https://travis-ci.org/jingweno/gh/builds/11911500", firstStatus.TargetURL)
}
