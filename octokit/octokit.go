package octokit

import (
	"net/http"
)

func NewClientWithToken(token string) *Client {
	return &Client{&http.Client{}, "", "", token}
}

func NewClientWithPassword(login, password string) *Client {
	return &Client{&http.Client{}, login, password, ""}
}

func NewClient() *Client {
	return &Client{&http.Client{}, "", "", ""}
}
