package octokit

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/bmizerany/assert"
)

func TestUsersService_GetCurrentUser(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("user.json"))
	})

	url, _ := CurrentUserURL.Expand(nil)
	user, result := client.Users(url).One()

	assert.T(t, !result.HasError())
	assert.Equal(t, 169064, user.ID)
	assert.Equal(t, "jingweno", user.Login)
	assert.Equal(t, "jingweno@gmail.com", user.Email)
	assert.Equal(t, "User", user.Type)
	assert.Equal(t, 17, user.Following)
	assert.Equal(t, 28, user.Followers)
	assert.Equal(t, 90, user.PublicRepos)
	assert.Equal(t, false, user.SiteAdmin)
	assert.Equal(t, "https://api.github.com/users/jingweno/repos", string(user.ReposURL))
}

func TestUsersService_UpdateCurrentUser(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testBody(t, r, "{\"email\":\"jingweno@gmail.com\"}\n")
		respondWithJSON(w, loadFixture("user.json"))
	})

	url, _ := CurrentUserURL.Expand(nil)
	userToUpdate := User{Email: "jingweno@gmail.com"}
	user, result := client.Users(url).Update(userToUpdate)

	assert.T(t, !result.HasError())
	assert.Equal(t, 169064, user.ID)
	assert.Equal(t, "jingweno", user.Login)
	assert.Equal(t, "jingweno@gmail.com", user.Email)
	assert.Equal(t, "User", user.Type)
	assert.Equal(t, "https://api.github.com/users/jingweno/repos", string(user.ReposURL))
}

func TestUsersService_GetUser(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/users/jingweno", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("user.json"))
	})

	url, err := UserURL.Expand(M{"user": "jingweno"})
	assert.Equal(t, nil, err)
	user, result := client.Users(url).One()

	assert.T(t, !result.HasError())
	assert.Equal(t, 169064, user.ID)
	assert.Equal(t, "jingweno", user.Login)
	assert.Equal(t, "jingweno@gmail.com", user.Email)
	assert.Equal(t, "User", user.Type)
	assert.Equal(t, "https://api.github.com/users/jingweno/repos", string(user.ReposURL))
}

func TestUsersService_All(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		rr := regexp.MustCompile(`users\?since=\d+`)
		assert.Tf(t, rr.MatchString(r.URL.String()), "Regexp should match users?since=\\d+")

		header := w.Header()
		link := fmt.Sprintf(`<%s>; rel="next", <%s>; rel="first"`, testURLOf("users?since=135"), testURLOf("users{?since}"))
		header.Set("Link", link)
		respondWithJSON(w, loadFixture("users.json"))
	})

	url, err := UserURL.Expand(M{"since": 1})
	assert.Equal(t, nil, err)

	q := url.Query()
	q.Set("since", "1")
	url.RawQuery = q.Encode()
	allUsers, result := client.Users(url).All()

	assert.T(t, !result.HasError())
	assert.Equal(t, 1, len(allUsers))
	assert.Equal(t, testURLStringOf("users?since=135"), string(*result.NextPage))

	nextPageURL, err := result.NextPage.Expand(nil)
	assert.Equal(t, nil, err)

	allUsers, result = client.Users(nextPageURL).All()
	assert.T(t, !result.HasError())
	assert.Equal(t, 1, len(allUsers))
}
