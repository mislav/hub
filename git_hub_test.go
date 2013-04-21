package main

import (
	"github.com/bmizerany/assert"
	"net/http"
	"os"
	"testing"
)

func TestCreatePullRequest(t *testing.T) {
	home := os.Getenv("HOME")
	config := loadConfig(home + "/.config/gh")

	client := &http.Client{}
	gh := GitHub{client, config.Token}
	params := PullRequestParams{"title", "body", "jingweno:master", "jingweno:pull_request"}
	err := gh.CreatePullRequest("jingweno", "gh", params)
	assert.Equal(t, nil, err)
}
