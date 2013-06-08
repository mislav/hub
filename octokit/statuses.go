package octokit

import (
	"fmt"
	"time"
)

type Status struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	State       string    `json:"state"`
	TargetUrl   string    `json:"target_url"`
	Description string    `json:"description"`
	Id          int       `json:"id"`
	Url         string    `json:"url"`
}

type StatusCreator struct {
	Login      string `json:"login"`
	Id         int    `json:"id"`
	AvatarUrl  string `json:"avatar_url"`
	GravatarId string `json:"gravatar_id"`
	Url        string `json:"url"`
}

func (c *Client) Statuses(repo Repository, sha string) ([]Status, error) {
	path := fmt.Sprintf("repos/%s/statuses/%s", repo, sha)
	var statuses []Status
	err := c.jsonGet(path, nil, &statuses)
	if err != nil {
		return nil, err
	}

	return statuses, nil
}
