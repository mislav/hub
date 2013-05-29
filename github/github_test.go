package github

import (
	"github.com/bmizerany/assert"
	"github.com/jingweno/gh/config"
	"net/http"
	"testing"
)

func _TestCreatePullRequest(t *testing.T) {
	config, _ := config.Load()

	client := &http.Client{}
	gh := GitHub{client, "jingweno", "123", config.Token}
	params := PullRequestParams{"title", "body", "jingweno:master", "jingweno:pull_request"}
	_, err := gh.CreatePullRequest(CurrentProject(), params)
	assert.Equal(t, nil, err)
}
