package octokit

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIssueCommentsService_AllIssueComments(t *testing.T) {
	setup()
	defer tearDown()

	link := fmt.Sprintf(`<%s>; rel="next", <%s>; rel="last"`, testURLOf("/repos/octokit/go-octokit/issues/1/comments?page=2"), testURLOf("/repos/octokit/go-octokit/issues/1/comments?page=3"))
	stubGet(t, "/repos/octokit/go-octokit/issues/1/comments", "issue_comments", map[string]string{"Link": link})

	comments, result := client.IssueComments().All(&IssueCommentsURL, M{"owner": "octokit", "repo": "go-octokit", "number": 1})
	assert.False(t, result.HasError())
	assert.Len(t, comments, 1)

	comment := comments[0]
	validateIssueComment(t, comment)

	assert.Equal(t, testURLStringOf("/repos/octokit/go-octokit/issues/1/comments?page=2"), string(*result.NextPage))
	assert.Equal(t, testURLStringOf("/repos/octokit/go-octokit/issues/1/comments?page=3"), string(*result.LastPage))

	validateNextPage_IssueComments(t, result)
}

func TestIssueCommentsService_AllRepoComments(t *testing.T) {
	setup()
	defer tearDown()

	link := fmt.Sprintf(`<%s>; rel="next", <%s>; rel="last"`, testURLOf("/repos/octokit/go-octokit/issues/comments?page=2"), testURLOf("/repos/octokit/go-octokit/issues/comments?page=3"))
	stubGet(t, "/repos/octokit/go-octokit/issues/comments", "issue_comments", map[string]string{"Link": link})

	comments, result := client.IssueComments().All(nil, M{"owner": "octokit", "repo": "go-octokit"})
	assert.False(t, result.HasError())
	assert.Len(t, comments, 1)

	comment := comments[0]
	validateIssueComment(t, comment)

	assert.Equal(t, testURLStringOf("/repos/octokit/go-octokit/issues/comments?page=2"), string(*result.NextPage))
	assert.Equal(t, testURLStringOf("/repos/octokit/go-octokit/issues/comments?page=3"), string(*result.LastPage))

	validateNextPage_IssueComments(t, result)
}

func TestIssueCommentsService_OneComment(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/repos/octokit/go-octokit/issues/comments/19158753", "issue_comment", nil)

	comment, result := client.IssueComments().One(nil, M{"owner": "octokit", "repo": "go-octokit", "id": 19158753})
	assert.False(t, result.HasError())

	validateIssueComment(t, *comment)
}

func TestIssueCommentsService_CreateComment(t *testing.T) {
	setup()
	defer tearDown()

	input := M{"body": "I am a comment"}
	wantReqBody, _ := json.Marshal(input)
	stubPost(t, "/repos/octokit/go-octokit/issues/1/comments", "issue_comment", nil, string(wantReqBody)+"\n", nil)

	comment, result := client.IssueComments().Create(nil, M{"owner": "octokit", "repo": "go-octokit", "number": 1}, input)
	assert.False(t, result.HasError())

	validateIssueComment(t, *comment)
}

func TestIssueCommentsService_UpdateComment(t *testing.T) {
	setup()
	defer tearDown()

	input := M{"body": "I am a comment"}
	wantReqBody, _ := json.Marshal(input)
	stubPatch(t, "/repos/octokit/go-octokit/issues/comments/19158753", "issue_comment", nil, string(wantReqBody)+"\n", nil)

	comment, result := client.IssueComments().Update(nil, M{"owner": "octokit", "repo": "go-octokit", "id": 19158753}, input)
	assert.False(t, result.HasError())

	validateIssueComment(t, *comment)
}

func TestIssueCommentsService_DeleteComment(t *testing.T) {
	setup()
	defer tearDown()

	var respHeaderParams = map[string]string{"Content-Type": "application/json"}
	stubDeletewCode(t, "/repos/octokit/go-octokit/issues/comments/19158753", respHeaderParams, 204)

	success, result := client.IssueComments().Delete(nil, M{"owner": "octokit", "repo": "go-octokit", "id": 19158753})
	assert.False(t, result.HasError())

	assert.True(t, success)
}

func TestIssueCommentsService_Failure(t *testing.T) {
	setup()
	defer tearDown()

	url := Hyperlink("}")
	comments, result := client.IssueComments().All(&url, nil)
	assert.True(t, result.HasError())
	assert.Len(t, comments, 0)

	comment, result := client.IssueComments().One(&url, nil)
	assert.True(t, result.HasError())
	assert.Nil(t, comment)

	comment, result = client.IssueComments().Create(&url, nil, nil)
	assert.True(t, result.HasError())
	assert.Nil(t, comment)

	comment, result = client.IssueComments().Update(&url, nil, nil)
	assert.True(t, result.HasError())
	assert.Nil(t, comment)

	success, result := client.IssueComments().Delete(&url, nil)
	assert.True(t, result.HasError())
	assert.False(t, success)
}

func validateIssueComment(t *testing.T, comment IssueComment) {
	testTime, _ := time.Parse("2006-01-02T15:04:05Z", "2013-06-09T00:53:41Z")

	assert.Equal(t, "https://api.github.com/repos/octokit/go-octokit/issues/comments/19158753", comment.URL)
	assert.Equal(t, "https://github.com/octokit/go-octokit/pull/1#issuecomment-19158753", comment.HTMLURL)
	assert.Equal(t, "https://api.github.com/repos/octokit/go-octokit/issues/1", comment.IssueURL)
	assert.Equal(t, 19158753, comment.ID)
	assert.Equal(t, &testTime, comment.CreatedAt)
	assert.Equal(t, &testTime, comment.UpdatedAt)
	assert.Equal(t, "Thanks! You're the first PR :). Merged", comment.Body)

	user := comment.User

	assert.Equal(t, "jingweno", user.Login)
	assert.Equal(t, 169064, user.ID)
}

func validateNextPage_IssueComments(t *testing.T, result *Result) {
	comments, result := client.IssueComments().All(result.NextPage, nil)
	assert.False(t, result.HasError())
	assert.Len(t, comments, 1)
}
