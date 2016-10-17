package octokit

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/bmizerany/assert"
)

func TestAuthorizationsService_One(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/authorizations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("authorization.json"))
	})

	url, err := AuthorizationsURL.Expand(M{"id": 1})
	assert.Equal(t, nil, err)

	auth, result := client.Authorizations(url).One()

	assert.T(t, !result.HasError())
	assert.Equal(t, 1, auth.ID)
	assert.Equal(t, "https://api.github.com/authorizations/1", auth.URL)
	assert.Equal(t, "456", auth.Token)
	assert.Equal(t, "", auth.Note)
	assert.Equal(t, "", auth.NoteURL)
	assert.Equal(t, "2012-11-16 01:05:51 +0000 UTC", auth.CreatedAt.String())
	assert.Equal(t, "2013-08-21 03:29:51 +0000 UTC", auth.UpdatedAt.String())

	app := App{ClientID: "123", URL: "http://localhost:8080", Name: "Test"}
	assert.Equal(t, app, auth.App)

	assert.Equal(t, 2, len(auth.Scopes))
	scopes := []string{"repo", "user"}
	assert.T(t, reflect.DeepEqual(auth.Scopes, scopes))
}

func TestAuthorizationsService_All(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/authorizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("authorizations.json"))
	})

	url, err := AuthorizationsURL.Expand(nil)
	assert.Equal(t, nil, err)

	auths, result := client.Authorizations(url).All()
	assert.T(t, !result.HasError())

	firstAuth := auths[0]
	assert.Equal(t, 1, firstAuth.ID)
	assert.Equal(t, "https://api.github.com/authorizations/1", firstAuth.URL)
	assert.Equal(t, "456", firstAuth.Token)
	assert.Equal(t, "", firstAuth.Note)
	assert.Equal(t, "", firstAuth.NoteURL)
	assert.Equal(t, "2012-11-16 01:05:51 +0000 UTC", firstAuth.CreatedAt.String())
	assert.Equal(t, "2013-08-21 03:29:51 +0000 UTC", firstAuth.UpdatedAt.String())

	app := App{ClientID: "123", URL: "http://localhost:8080", Name: "Test"}
	assert.Equal(t, app, firstAuth.App)

	assert.Equal(t, 2, len(firstAuth.Scopes))
	scopes := []string{"repo", "user"}
	assert.T(t, reflect.DeepEqual(firstAuth.Scopes, scopes))
}

func TestAuthorizationsService_Create(t *testing.T) {
	setup()
	defer tearDown()

	params := AuthorizationParams{Scopes: []string{"public_repo"}}

	mux.HandleFunc("/authorizations", func(w http.ResponseWriter, r *http.Request) {
		var authParams AuthorizationParams
		json.NewDecoder(r.Body).Decode(&authParams)
		assert.T(t, reflect.DeepEqual(authParams, params))

		testMethod(t, r, "POST")
		respondWithJSON(w, loadFixture("create_authorization.json"))
	})

	url, err := AuthorizationsURL.Expand(nil)
	assert.Equal(t, nil, err)

	auth, _ := client.Authorizations(url).Create(params)

	assert.Equal(t, 3844190, auth.ID)
	assert.Equal(t, "https://api.github.com/authorizations/3844190", auth.URL)
	assert.Equal(t, "123", auth.Token)
	assert.Equal(t, "", auth.Note)
	assert.Equal(t, "", auth.NoteURL)
	assert.Equal(t, "2013-09-28 18:44:39 +0000 UTC", auth.CreatedAt.String())
	assert.Equal(t, "2013-09-28 18:44:39 +0000 UTC", auth.UpdatedAt.String())

	app := App{ClientID: "00000000000000000000", URL: "http://developer.github.com/v3/oauth/#oauth-authorizations-api", Name: "GitHub API"}
	assert.Equal(t, app, auth.App)

	assert.Equal(t, 1, len(auth.Scopes))
	scopes := []string{"public_repo"}
	assert.T(t, reflect.DeepEqual(auth.Scopes, scopes))
}
