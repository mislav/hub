package octokit

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthorizationsService_One(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/authorizations/1", "authorization", nil)

	url, err := AuthorizationsURL.Expand(M{"id": 1})
	assert.NoError(t, err)

	auth, result := client.Authorizations(url).One()

	assert.False(t, result.HasError())
	assert.Equal(t, 1, auth.ID)
	assert.Equal(t, "https://api.github.com/authorizations/1", auth.URL)
	assert.Equal(t, "456", auth.Token)
	assert.Equal(t, "", auth.Note)
	assert.Equal(t, "", auth.NoteURL)
	assert.Equal(t, "2012-11-16 01:05:51 +0000 UTC", auth.CreatedAt.String())
	assert.Equal(t, "2013-08-21 03:29:51 +0000 UTC", auth.UpdatedAt.String())

	app := App{ClientID: "123", URL: "http://localhost:8080", Name: "Test"}
	assert.Equal(t, app, auth.App)

	assert.Len(t, auth.Scopes, 2)
	scopes := []string{"repo", "user"}
	assert.True(t, reflect.DeepEqual(auth.Scopes, scopes))
}

func TestAuthorizationsService_All(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/authorizations", "authorizations", nil)

	url, err := AuthorizationsURL.Expand(nil)
	assert.NoError(t, err)

	auths, result := client.Authorizations(url).All()
	assert.False(t, result.HasError())

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

	assert.Len(t, firstAuth.Scopes, 2)
	scopes := []string{"repo", "user"}
	assert.True(t, reflect.DeepEqual(firstAuth.Scopes, scopes))
}

func TestAuthorizationsService_Create(t *testing.T) {
	setup()
	defer tearDown()

	params := AuthorizationParams{Scopes: []string{"public_repo"}}

	wantReqBody, _ := json.Marshal(params)
	stubPost(t, "/authorizations", "create_authorization", nil, string(wantReqBody)+"\n", nil)

	url, err := AuthorizationsURL.Expand(nil)
	assert.NoError(t, err)

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

	assert.Len(t, auth.Scopes, 1)
	scopes := []string{"public_repo"}
	assert.True(t, reflect.DeepEqual(auth.Scopes, scopes))
}
