package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GitHubError struct {
	Resource string `json:"resource"`
	Field    string `json:"field"`
	Value    string `json:"value"`
	Code     string `json:"code"`
	Message  string `json:"message"`
}

type GitHubErrors struct {
	Message string        `json:"message"`
	Errors  []GitHubError `json:"errors"`
}

func performBasicAuth(gh *GitHub) error {
	user := gh.config.FetchUser()
	password := gh.config.FetchPassword()
	gh.updateBasicAuth(user, password)

	return obtainOAuthTokenWithBasicAuth(gh)
}

func obtainOAuthTokenWithBasicAuth(gh *GitHub) error {
	auths, err := listAuthorizations(gh)
	if err != nil {
		return err
	}

	var token string
	for _, auth := range auths {
		if auth.NoteUrl == OAuthAppUrl {
			token = auth.Token
			break
		}
	}

	if token == "" {
		authParam := AuthorizationParams{}
		authParam.Scopes = append(authParam.Scopes, "repo")
		authParam.Note = "gh"
		authParam.NoteUrl = OAuthAppUrl

		auth, err := createAuthorization(gh, authParam)
		if err != nil {
			return err
		}

		token = auth.Token
	}

	gh.updateToken(token)

	return nil
}

func httpGet(gh *GitHub, uri string, extraHeaders map[string]string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", GitHubApiUrl, uri)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if extraHeaders != nil {
		for h, v := range extraHeaders {
			request.Header.Set(h, v)
		}
	}

	return performRequest(gh, request)
}

func httpPost(gh *GitHub, uri string, extraHeaders map[string]string, content *bytes.Buffer) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", GitHubApiUrl, uri)
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

	return performRequest(gh, request)
}

func performRequest(gh *GitHub, request *http.Request) (*http.Response, error) {
	if gh.authorization == "" {
		err := performBasicAuth(gh)
		if err != nil {
			return nil, err
		}
	}

	request.Header.Set("Authorization", gh.authorization)

	response, err := gh.httpClient.Do(request)
	if err != nil {
		return response, err
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return response, err
	}

	err = handleGitHubErrors(response)

	return nil, err
}

func handleGitHubErrors(response *http.Response) error {
	var githubErrors GitHubErrors
	err := unmarshalBody(response, &githubErrors)
	if err != nil {
		return err
	}

	errorMessages := make([]string, len(githubErrors.Errors))
	for _, e := range githubErrors.Errors {
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

	var text string
	for _, m := range errorMessages {
		if m != "" {
			text = text + m + "\n"
		}
	}

	if text == "" {
		text = githubErrors.Message
	}

	return errors.New(text)
}

func unmarshalBody(response *http.Response, v interface{}) error {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return unmarshal(body, v)
}

func unmarshal(body []byte, v interface{}) error {
	return json.Unmarshal(body, v)
}
