package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const (
	GitHubUrl string = "https://api.github.com"
)

type unprocessableEntityError struct {
	Resource string `json:"resource"`
	Field    string `json:"field"`
	Value    string `json:"value"`
	Code     string `json:"code"`
	Message  string `json:"message"`
}

type unprocessableEntity struct {
	Message string                     `json:"message"`
	Errors  []unprocessableEntityError `json:"errors"`
}

func NewGitHub() *GitHub {
	configFile := filepath.Join(os.Getenv("HOME"), ".config", "gh")
	config, _ := LoadConfig(configFile)

	return &GitHub{&http.Client{}, config.Token}
}

type GitHub struct {
	httpClient    *http.Client
	Authorization string
}

func (gh *GitHub) performBasicAuth(url *url.URL) {
	url.String()
}

func (gh *GitHub) call(request *http.Request) (*http.Response, error) {
	if len(gh.Authorization) == 0 {
		gh.performBasicAuth(request.URL)
	}

	request.Header.Set("Authorization", "token "+gh.Authorization)

	response, err := gh.httpClient.Do(request)
	if err != nil {
		return response, err
	}

	if response.StatusCode != 422 {
		return response, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		var unprocessable unprocessableEntity
		err = json.Unmarshal(body, &unprocessable)
		if err != nil {
			return response, err
		}

		errorMessages := make([]string, len(unprocessable.Errors))
		for _, e := range unprocessable.Errors {
			switch e.Code {
			case "custom":
				errorMessages = append(errorMessages, e.Message)
			case "missing_field":
				errorMessages = append(errorMessages, "Missing field: "+e.Field)
			case "invalid":
				errorMessages = append(errorMessages, "Invalid value for "+e.Field+": "+e.Value)
			case "unauthorized":
				errorMessages = append(errorMessages, "Not allow to change field "+e.Field)
			}
		}

		text := ""
		for _, m := range errorMessages {
			if len(m) > 0 {
				text = text + m + "\n"
			}
		}

		if len(text) == 0 {
			text = unprocessable.Message
		}

		err = errors.New(text)
	}

	return response, err
}

func (gh *GitHub) httpPost(uri string, extraHeaders map[string]string, content *bytes.Buffer) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", GitHubUrl, uri)
	request, err := http.NewRequest("POST", url, content)
	if err != nil {
		return nil, err
	}

	if extraHeaders != nil {
		for h, v := range extraHeaders {
			request.Header.Set(h, v)
		}
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	return gh.call(request)
}

type PullRequestParams struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Base  string `json:"base"`
	Head  string `json:"head"`
}

type PullRequestResponse struct {
	Url      string `json:"url"`
	HtmlUrl  string `json:"html_url"`
	DiffUrl  string `json:"diff_url"`
	PatchUrl string `json:"patch_url"`
	IssueUrl string `json:"issue_url"`
}

func (gh *GitHub) CreatePullRequest(owner, repo string, params PullRequestParams) (*PullRequestResponse, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(b)
	url := "/repos/" + owner + "/" + repo + "/pulls"
	response, err := gh.httpPost(url, nil, buffer)
	if err != nil {
		return nil, err
	}

	js, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var pullRequestResponse PullRequestResponse
	err = json.Unmarshal(js, &pullRequestResponse)
	if err != nil {
		return nil, err
	}

	return &pullRequestResponse, nil
}
