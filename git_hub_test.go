package main

import (
	"github.com/bmizerany/assert"
	"net/http"
	"testing"
)

func _TestCreatePullRequest(t *testing.T) {
	config, _ := LoadConfig("./test_support/gh")

	client := &http.Client{}
	gh := GitHub{client, "jingweno", "123", config.Token}
	params := PullRequestParams{"title", "body", "jingweno:master", "jingweno:pull_request"}
	_, err := gh.CreatePullRequest("jingweno", "gh", params)
	assert.Equal(t, nil, err)
}
