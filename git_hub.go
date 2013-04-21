package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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

type GitHub struct {
	httpClient    *http.Client
	Authorization string
}

func (gh *GitHub) call(request *http.Request) (response *http.Response, err error) {
	request.Header.Set("Authorization", "token "+gh.Authorization)

	response, err = gh.httpClient.Do(request)
	if err != nil {
		return
	}

	if response.StatusCode != 422 {
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		var unprocessable unprocessableEntity
		err = json.Unmarshal(body, &unprocessable)
		if err != nil {
			return
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

	return
}

func (gh *GitHub) httpPost(uri string, extraHeaders map[string]string, content *bytes.Buffer) (response *http.Response, err error) {
	url := fmt.Sprintf("%s%s", GitHubUrl, uri)
	request, err := http.NewRequest("POST", url, content)
	if err != nil {
		return
	}

	// Add (any of) the extra headers to the request
	if extraHeaders != nil {
		for h, v := range extraHeaders {
			request.Header.Set(h, v)
		}
	}

	// Set the Content-Type header
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	response, err = gh.call(request)

	return
}

type PullRequestParams struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Base  string `json:"base"`
	Head  string `json:"head"`
}

func (gh *GitHub) CreatePullRequest(owner, repo string, params PullRequestParams) (err error) {
	b, err := json.Marshal(params)
	if err != nil {
		return
	}

	buffer := bytes.NewBuffer(b)
	url := "/repos/" + owner + "/" + repo + "/pulls"
	_, err = gh.httpPost(url, nil, buffer)

	return err
}
