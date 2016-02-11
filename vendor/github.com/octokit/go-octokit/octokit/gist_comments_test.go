package octokit

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGistCommentsService_AllComments(t *testing.T) {
	setup()
	defer tearDown()

	link := fmt.Sprintf(`<%s>; rel="next", <%s>; rel="last"`, testURLOf("/gists/1721489/comments?page=2"), testURLOf("/gists/1721489/comments?page=3"))
	stubGet(t, "/gists/1721489/comments", "gist_comments", map[string]string{"Link": link})

	comments, result := client.GistComments().All(nil, M{"gist_id": 1721489})
	fmt.Println(result.Error())
	assert.False(t, result.HasError())
	assert.Len(t, comments, 1)

	comment := comments[0]
	validateGistComment(t, comment)

	assert.Equal(t, testURLStringOf("/gists/1721489/comments?page=2"), string(*result.NextPage))
	assert.Equal(t, testURLStringOf("/gists/1721489/comments?page=3"), string(*result.LastPage))

	validateNextPage_GistComments(t, result)
}

func TestGistCommentsService_OneComment(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/gists/1721489/comments/1199157", "gist_comment", nil)

	comment, result := client.GistComments().One(nil, M{"gist_id": 1721489, "id": 1199157})
	assert.False(t, result.HasError())

	validateGistComment(t, *comment)
}

func TestGistCommentsService_CreateComment(t *testing.T) {
	setup()
	defer tearDown()

	input := M{"body": "I am a comment"}
	wantReqBody, _ := json.Marshal(input)
	stubPost(t, "/gists/1721489/comments", "gist_comment", nil, string(wantReqBody)+"\n", nil)

	comment, result := client.GistComments().Create(nil, M{"gist_id": 1721489}, input)
	assert.False(t, result.HasError())

	validateGistComment(t, *comment)
}

func TestGistCommentsService_UpdateComment(t *testing.T) {
	setup()
	defer tearDown()

	input := M{"body": "I am a comment"}
	wantReqBody, _ := json.Marshal(input)
	stubPatch(t, "/gists/1721489/comments/1199157", "gist_comment", nil, string(wantReqBody)+"\n", nil)

	comment, result := client.GistComments().Update(nil, M{"gist_id": 1721489, "id": 1199157}, input)
	assert.False(t, result.HasError())

	validateGistComment(t, *comment)
}

func TestGistCommentsService_DeleteComment(t *testing.T) {
	setup()
	defer tearDown()

	respHeaderParams := map[string]string{"Content-Type": "application/json"}
	stubDeletewCode(t, "/gists/1721489/comments/1199157", respHeaderParams, 204)

	success, result := client.GistComments().Delete(nil, M{"gist_id": 1721489, "id": 1199157})
	assert.False(t, result.HasError())

	assert.True(t, success)
}

func TestGistCommentsService_Failure(t *testing.T) {
	setup()
	defer tearDown()

	url := Hyperlink("}")
	comments, result := client.GistComments().All(&url, nil)
	assert.True(t, result.HasError())
	assert.Len(t, comments, 0)

	comment, result := client.GistComments().One(&url, nil)
	assert.True(t, result.HasError())
	assert.Nil(t, comment)

	comment, result = client.GistComments().Create(&url, nil, nil)
	assert.True(t, result.HasError())
	assert.Nil(t, comment)

	comment, result = client.GistComments().Update(&url, nil, nil)
	assert.True(t, result.HasError())
	assert.Nil(t, comment)

	success, result := client.GistComments().Delete(&url, nil)
	assert.True(t, result.HasError())
	assert.False(t, success)
}

func validateGistComment(t *testing.T, comment GistComment) {
	testTime, _ := time.Parse("2006-01-02T15:04:05Z", "2014-03-26T07:30:04Z")

	assert.Equal(t, "https://api.github.com/gists/1721489/comments/1199157", comment.URL)
	assert.Equal(t, 1199157, comment.ID)
	assert.Equal(t, &testTime, comment.CreatedAt)
	assert.Equal(t, &testTime, comment.UpdatedAt)
	assert.Equal(t, "This is a body", comment.Body)

	user := comment.User

	assert.Equal(t, "purse9644", user.Login)
	assert.Equal(t, 7067053, user.ID)
}

func validateNextPage_GistComments(t *testing.T, result *Result) {
	comments, result := client.GistComments().All(result.NextPage, nil)
	assert.False(t, result.HasError())
	assert.Len(t, comments, 1)
}
