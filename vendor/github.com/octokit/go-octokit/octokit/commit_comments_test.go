package octokit

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCommitCommentsService_AllRepoComments(t *testing.T) {
	setup()
	defer tearDown()

	link := fmt.Sprintf(`<%s>; rel="next", <%s>; rel="last"`, testURLOf("/repos/octokit/go-octokit/comments?page=2"), testURLOf("/repos/octokit/go-octokit/comments?page=3"))
	stubGet(t, "/repos/octokit/go-octokit/comments", "commit_comments", map[string]string{"Link": link})

	comments, result := client.CommitComments().All(nil, M{"owner": "octokit", "repo": "go-octokit"})
	assert.False(t, result.HasError())
	assert.Len(t, comments, 1)

	comment := comments[0]
	validateCommitComment(t, comment)

	assert.Equal(t, testURLStringOf("/repos/octokit/go-octokit/comments?page=2"), string(*result.NextPage))
	assert.Equal(t, testURLStringOf("/repos/octokit/go-octokit/comments?page=3"), string(*result.LastPage))

	validateNextPage_CommitComments(t, result)
}

func TestCommitCommentsService_AllCommitComments(t *testing.T) {
	setup()
	defer tearDown()

	link := fmt.Sprintf(`<%s>; rel="next", <%s>; rel="last"`, testURLOf("/repos/octokit/go-octokit/commits/8b8347dc11c81b64fdd9938d34dc4ef6a07dbf09/comments?page=2"), testURLOf("/repos/octokit/go-octokit/commits/8b8347dc11c81b64fdd9938d34dc4ef6a07dbf09/comments?page=3"))
	stubGet(t, "/repos/octokit/go-octokit/commits/8b8347dc11c81b64fdd9938d34dc4ef6a07dbf09/comments", "commit_comments", map[string]string{"Link": link})

	comments, result := client.CommitComments().All(&CommitCommentsURL, M{"owner": "octokit", "repo": "go-octokit", "sha": "8b8347dc11c81b64fdd9938d34dc4ef6a07dbf09"})
	assert.False(t, result.HasError())
	assert.Len(t, comments, 1)

	comment := comments[0]
	validateCommitComment(t, comment)

	assert.Equal(t, testURLStringOf("/repos/octokit/go-octokit/commits/8b8347dc11c81b64fdd9938d34dc4ef6a07dbf09/comments?page=2"), string(*result.NextPage))
	assert.Equal(t, testURLStringOf("/repos/octokit/go-octokit/commits/8b8347dc11c81b64fdd9938d34dc4ef6a07dbf09/comments?page=3"), string(*result.LastPage))

	validateNextPage_CommitComments(t, result)
}

func TestCommitCommentsService_OneComment(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/repos/octokit/go-octokit/comments/4236029", "commit_comment", nil)

	comment, result := client.CommitComments().One(nil, M{"owner": "octokit", "repo": "go-octokit", "id": 4236029})
	assert.False(t, result.HasError())

	validateCommitComment(t, *comment)
}

func TestCommitCommentsService_CreateComment(t *testing.T) {
	setup()
	defer tearDown()

	input := M{
		"body":     "I am a comment",
		"path":     "root.go",
		"position": 46,
	}

	wantReqBody, _ := json.Marshal(input)
	stubPost(t, "/repos/octokit/go-octokit/commits/8b8347dc11c81b64fdd9938d34dc4ef6a07dbf09/comments",
		"commit_comment", nil, string(wantReqBody)+"\n", nil)

	comment, result := client.CommitComments().Create(nil, M{"owner": "octokit", "repo": "go-octokit", "sha": "8b8347dc11c81b64fdd9938d34dc4ef6a07dbf09"}, input)
	assert.False(t, result.HasError())

	validateCommitComment(t, *comment)
}

func TestCommitCommentsService_UpdateComment(t *testing.T) {
	setup()
	defer tearDown()

	input := M{"body": "I am a comment"}
	wantReqBody, _ := json.Marshal(input)
	stubPatch(t, "/repos/octokit/go-octokit/comments/4236029", "commit_comment", nil, string(wantReqBody)+"\n", nil)

	comment, result := client.CommitComments().Update(nil, M{"owner": "octokit", "repo": "go-octokit", "id": 4236029}, input)
	assert.False(t, result.HasError())

	validateCommitComment(t, *comment)
}

func TestCommitCommentsService_DeleteComment(t *testing.T) {
	setup()
	defer tearDown()

	var respHeaderParams = map[string]string{"Content-Type": "application/json"}
	stubDeletewCode(t, "/repos/octokit/go-octokit/comments/4236029", respHeaderParams, 204)

	success, result := client.CommitComments().Delete(nil, M{"owner": "octokit", "repo": "go-octokit", "id": 4236029})
	assert.False(t, result.HasError())

	assert.True(t, success)
}

func TestCommitCommentsService_Failure(t *testing.T) {
	setup()
	defer tearDown()

	url := Hyperlink("}")
	comments, result := client.CommitComments().All(&url, nil)
	assert.True(t, result.HasError())
	assert.Len(t, comments, 0)

	comment, result := client.CommitComments().One(&url, nil)
	assert.True(t, result.HasError())
	assert.Nil(t, comment)

	comment, result = client.CommitComments().Create(&url, nil, nil)
	assert.True(t, result.HasError())
	assert.Nil(t, comment)

	comment, result = client.CommitComments().Update(&url, nil, nil)
	assert.True(t, result.HasError())
	assert.Nil(t, comment)

	success, result := client.CommitComments().Delete(&url, nil)
	assert.True(t, result.HasError())
	assert.False(t, success)
}

func validateCommitComment(t *testing.T, comment CommitComment) {
	testTime, _ := time.Parse("2006-01-02T15:04:05Z", "2013-10-02T19:32:40Z")

	assert.Equal(t, "https://api.github.com/repos/octokit/go-octokit/comments/4236029", comment.URL)
	assert.Equal(t, "https://github.com/octokit/go-octokit/commit/8b8347dc11c81b64fdd9938d34dc4ef6a07dbf09#commitcomment-4236029", comment.HTMLURL)
	assert.Equal(t, 4236029, comment.ID)
	assert.Equal(t, 46, comment.Position)
	assert.Equal(t, 46, comment.Line)
	assert.Equal(t, "root.go", comment.Path)
	assert.Equal(t, "8b8347dc11c81b64fdd9938d34dc4ef6a07dbf09", comment.CommitID)
	assert.Equal(t, &testTime, comment.CreatedAt)
	assert.Equal(t, &testTime, comment.UpdatedAt)
	assert.Equal(t, ":heart:\r\n\r\nAre you handling plain `url`, too? In Octokit.rb, we parse those as a `self` relation.", comment.Body)

	user := comment.User

	assert.Equal(t, "pengwynn", user.Login)
	assert.Equal(t, 865, user.ID)
}

func validateNextPage_CommitComments(t *testing.T, result *Result) {
	comments, result := client.CommitComments().All(result.NextPage, nil)
	assert.False(t, result.HasError())
	assert.Len(t, comments, 1)
}
