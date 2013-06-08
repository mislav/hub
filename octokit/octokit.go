package octokit

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	GitHubHost    string = "github.com"
	GitHubApiUrl  string = "https://" + GitHubApiHost
	GitHubApiHost string = "api.github.com"
	OAuthAppUrl   string = "http://owenou.com/gh"
)

type GitHubError struct {
	Resource string      `json:"resource"`
	Field    string      `json:"field"`
	Value    interface{} `json:"value"`
	Code     string      `json:"code"`
	Message  string      `json:"message"`
}

type GitHubErrors struct {
	Message string        `json:"message"`
	Errors  []GitHubError `json:"errors"`
}

type Client struct {
	httpClient *http.Client
	Login      string
	Password   string
	Token      string
}

func (c *Client) get(path string, extraHeaders map[string]string) ([]byte, error) {
	return c.request("GET", path, extraHeaders, nil)
}

func (c *Client) post(path string, extraHeaders map[string]string, content *bytes.Buffer) ([]byte, error) {
	return c.request("POST", path, extraHeaders, content)
}

func (c *Client) request(method, path string, extraHeaders map[string]string, content io.Reader) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", GitHubApiUrl, path)
	request, err := http.NewRequest(method, url, content)
	if err != nil {
		return nil, err
	}

	c.setDefaultHeaders(request)

	if extraHeaders != nil {
		for h, v := range extraHeaders {
			request.Header.Set(h, v)
		}
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 && response.StatusCode < 600 {
		return nil, handleErrors(body)
	}

	return body, nil
}

func (c *Client) setDefaultHeaders(request *http.Request) {
	request.Header.Set("Accept", "application/vnd.github.beta+json")
	if c.Login != "" && c.Password != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Basic %s", hashAuth(c.Login, c.Password)))
	}
	if c.Token != "" {
		request.Header.Set("Authorization", fmt.Sprintf("token %s", c.Token))
	}
}

func handleErrors(body []byte) error {
	var githubErrors GitHubErrors
	err := json.Unmarshal(body, &githubErrors)
	if err != nil {
		return err
	}

	errorMessages := make([]string, len(githubErrors.Errors))
	for _, e := range githubErrors.Errors {
		switch e.Code {
		case "custom":
			errorMessages = append(errorMessages, e.Message)
		case "missing_field":
			msg := fmt.Sprintf("Missing field: %s", e.Field)
			errorMessages = append(errorMessages, msg)
		case "invalid":
			msg := fmt.Sprintf("Invalid value for %s: %v", e.Field, e.Value)
			errorMessages = append(errorMessages, msg)
		case "unauthorized":
			errorMessages = append(errorMessages, "Not allow to change field "+e.Field)
		}
	}

	text := strings.Join(errorMessages, "\n")
	if text == "" {
		text = githubErrors.Message
	}

	return errors.New(text)
}

func hashAuth(u, p string) string {
	var a = fmt.Sprintf("%s:%s", u, p)
	return base64.StdEncoding.EncodeToString([]byte(a))
}

func NewClientWithPassword(login, password string) *Client {
	return &Client{&http.Client{}, login, password, ""}
}

func NewClient() *Client {
	return &Client{&http.Client{}, "", "", ""}
}
