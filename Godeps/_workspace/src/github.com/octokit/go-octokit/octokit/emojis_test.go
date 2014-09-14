package octokit

import (
	"net/http"
	"testing"

	"github.com/bmizerany/assert"
)

func TestRootEmojisService_All(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/emojis", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("emojis.json"))
	})

	url, err := EmojisURL.Expand(nil)
	assert.Equal(t, nil, err)

	emojis, result := client.Emojis(url).All()
	assert.T(t, !result.HasError())
	var penguin = "https://github.global.ssl.fastly.net/images/icons/emoji/penguin.png?v5"
	var metal = "https://github.global.ssl.fastly.net/images/icons/emoji/metal.png?v5"
	assert.Equal(t, penguin, emojis["penguin"])
	assert.Equal(t, metal, emojis["metal"])
}
