package octokit

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIssueLabelsService_Add(t *testing.T) {
	setup()
	defer tearDown()

	input := []string{"newLabel", "anotherNewLabel"}
	wantReqBody, _ := json.Marshal(input)
	stubPost(t, "/repos/octokit/go-octokit/issues/33/labels", "issue_labels_added", nil, string(wantReqBody)+"\n", nil)

	labels, result := client.IssueLabels().Add(nil, M{"owner": "octokit", "repo": "go-octokit", "number": 33}, input)

	assert.False(t, result.HasError())

	assert.Equal(t, 2, len(labels))

	assert.Equal(t, "https://api.github.com/repos/octokit/go-octokit/labels/newLabel", labels[0].URL)
	assert.Equal(t, "newLabel", labels[0].Name)
	assert.Equal(t, "ffffff", labels[0].Color)

	assert.Equal(t, "https://api.github.com/repos/octokit/go-octokit/labels/anotherNewLabel", labels[1].URL)
	assert.Equal(t, "anotherNewLabel", labels[1].Name)
	assert.Equal(t, "000000", labels[1].Color)
}

func TestIssueLabelsService_All(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/repos/octokit/go-octokit/issues/33/labels", "issue_labels", nil)

	labels, result := client.IssueLabels().All(nil, M{"owner": "octokit", "repo": "go-octokit", "number": 33})

	assert.False(t, result.HasError())

	assert.Equal(t, 2, len(labels))

	assert.Equal(t, "https://api.github.com/repos/octokit/go-octokit/labels/bug", labels[0].URL)
	assert.Equal(t, "bug", labels[0].Name)
	assert.Equal(t, "fc2929", labels[0].Color)

	assert.Equal(t, "https://api.github.com/repos/octokit/go-octokit/labels/duplicate", labels[1].URL)
	assert.Equal(t, "duplicate", labels[1].Name)
	assert.Equal(t, "cccccc", labels[1].Color)
}

func TestIssueLabelsService_Remove(t *testing.T) {
	setup()
	defer tearDown()

	var respHeaderParams = map[string]string{"Content-Type": "application/json"}
	stubDeletewCode(t, "/repos/octokit/go-octokit/issues/33/labels/theName", respHeaderParams, 204)

	success, result := client.IssueLabels().Remove(nil, M{"owner": "octokit", "repo": "go-octokit", "number": 33, "name": "theName"})

	assert.False(t, result.HasError())

	assert.True(t, success)
}

func TestIssueLabelsService_RemoveAll(t *testing.T) {
	setup()
	defer tearDown()

	var respHeaderParams = map[string]string{"Content-Type": "application/json"}
	stubDeletewCode(t, "/repos/octokit/go-octokit/issues/33/labels", respHeaderParams, 204)

	success, result := client.IssueLabels().RemoveAll(nil, M{"owner": "octokit", "repo": "go-octokit", "number": 33})

	assert.False(t, result.HasError())

	assert.True(t, success)
}

func TestIssueLabelsService_ReplaceAll(t *testing.T) {
	setup()
	defer tearDown()

	respHeaderParams := map[string]string{"Content-Type": "application/json"}
	stubPutwCode(t, "/repos/octokit/go-octokit/issues/33/labels", "issue_labels_replaced", nil, "", respHeaderParams, 200)

	labels, result := client.IssueLabels().ReplaceAll(nil, M{"owner": "octokit", "repo": "go-octokit", "number": 33}, []string{"theFirstNewLabel", "theSecondNewLabel"})

	assert.False(t, result.HasError())

	assert.Equal(t, 2, len(labels))

	assert.Equal(t, "https://api.github.com/repos/octokit/go-octokit/labels/theFirstNewLabel", labels[0].URL)
	assert.Equal(t, "theFirstNewLabel", labels[0].Name)
	assert.Equal(t, "000000", labels[0].Color)

	assert.Equal(t, "https://api.github.com/repos/octokit/go-octokit/labels/theSecondNewLabel", labels[1].URL)
	assert.Equal(t, "theSecondNewLabel", labels[1].Name)
	assert.Equal(t, "ffffff", labels[1].Color)
}
