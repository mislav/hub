package octokit

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLabelsService_All(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/repos/octokit/go-octokit/labels", "labels", nil)

	labels, result := client.Labels().All(nil, M{"owner": "octokit", "repo": "go-octokit"})

	assert.False(t, result.HasError())

	assert.Equal(t, 2, len(labels))

	assert.Equal(t, "https://api.github.com/repos/octokit/go-octokit/labels/bug", labels[0].URL)
	assert.Equal(t, "bug", labels[0].Name)
	assert.Equal(t, "fc2929", labels[0].Color)

	assert.Equal(t, "https://api.github.com/repos/octokit/go-octokit/labels/duplicate", labels[1].URL)
	assert.Equal(t, "duplicate", labels[1].Name)
	assert.Equal(t, "cccccc", labels[1].Color)
}

func TestLabelsService_One(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/repos/octokit/go-octokit/labels/bug", "label", nil)

	label, result := client.Labels().One(nil, M{"owner": "octokit", "repo": "go-octokit", "name": "bug"})

	assert.False(t, result.HasError())

	assert.Equal(t, "https://api.github.com/repos/octokit/go-octokit/labels/bug", (*label).URL)
	assert.Equal(t, "bug", (*label).Name)
	assert.Equal(t, "fc2929", (*label).Color)
}

func TestLabelsService_Create(t *testing.T) {
	setup()
	defer tearDown()

	input := M{"name": "theLabel", "color": "ffffff"}
	wantReqBody, _ := json.Marshal(input)
	stubPost(t, "/repos/octokit/go-octokit/labels", "label_created", nil, string(wantReqBody)+"\n", nil)

	label, result := client.Labels().Create(nil, M{"owner": "octokit", "repo": "go-octokit"}, input)

	assert.False(t, result.HasError())

	assert.Equal(t, "https://api.github.com/repos/octokit/go-octokit/labels/theLabel", label.URL)
	assert.Equal(t, "theLabel", label.Name)
	assert.Equal(t, "ffffff", label.Color)
}

func TestLabelsService_Update(t *testing.T) {
	setup()
	defer tearDown()

	input := M{"name": "theChangedName", "color": "000000"}
	wantReqBody, _ := json.Marshal(input)
	stubPatch(t, "/repos/octokit/go-octokit/labels/theName", "label_changed", nil, string(wantReqBody)+"\n", nil)

	label, result := client.Labels().Update(nil, M{"owner": "octokit", "repo": "go-octokit", "name": "theName"}, input)

	assert.False(t, result.HasError())

	assert.Equal(t, "https://api.github.com/repos/octokit/go-octokit/labels/theChangedName", label.URL)
	assert.Equal(t, "theChangedName", label.Name)
	assert.Equal(t, "000000", label.Color)
}

func TestLabelsService_Delete(t *testing.T) {
	setup()
	defer tearDown()

	var respHeaderParams = map[string]string{"Content-Type": "application/json"}
	stubDeletewCode(t, "/repos/octokit/go-octokit/labels/theLabel", respHeaderParams, 204)

	success, result := client.Labels().Delete(nil, M{"owner": "octokit", "repo": "go-octokit", "name": "theLabel"})

	assert.False(t, result.HasError())

	assert.True(t, success)
}
