package github

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/jingweno/gh/config"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	GitHubUrl   string = "https://" + GitHubHost
	GitHubHost  string = "api.github.com"
	OAuthAppUrl string = "http://owenou.com/gh"
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

type App struct {
	Url      string `json:"url"`
	Name     string `json:"name"`
	ClientId string `json:"client_id"`
}

type Authorization struct {
	Scopes  []string `json:"scopes"`
	Url     string   `json:"url"`
	App     App      `json:"app"`
	Token   string   `json:"token"`
	Note    string   `josn:"note"`
	NoteUrl string   `josn:"note_url"`
}

func New() *GitHub {
	project := CurrentProject()
	config, err := config.Load(project.Owner)

	var user, auth string
	if err == nil {
		user = config.User
		auth = config.Token
	}

	if auth != "" {
		auth = fmt.Sprintf("token %s", auth)
	}

	return &GitHub{&http.Client{}, user, "", auth, project}
}

func hashAuth(u, p string) string {
	var a = fmt.Sprintf("%s:%s", u, p)
	return base64.StdEncoding.EncodeToString([]byte(a))
}

type GitHub struct {
	httpClient    *http.Client
	User          string
	Password      string
	Authorization string
	Project       *Project
}

func (gh *GitHub) performBasicAuth() error {
	user := gh.User
	if user == "" {
		user = CurrentProject().Owner
		gh.User = user
	}
	if user == "" {
		// TODO: prompt user
		log.Fatal("TODO: prompt user for basic auth")
	}

	msg := fmt.Sprintf("%s password for %s (never stored): ", GitHubHost, user)
	fmt.Print(msg)

	pass := gopass.GetPasswd()
	if len(pass) == 0 {
		return errors.New("Password cannot be empty.")
	}
	gh.Password = string(pass)

	return gh.obtainOAuthTokenWithBasicAuth()
}

func (gh *GitHub) obtainOAuthTokenWithBasicAuth() error {
	gh.Authorization = fmt.Sprintf("Basic %s", hashAuth(gh.User, gh.Password))
	response, err := gh.httpGet("/authorizations", nil)
	if err != nil {
		return err
	}

	var auths []Authorization
	err = unmarshalBody(response, &auths)
	if err != nil {
		return err
	}

	var token string
	for _, auth := range auths {
		if auth.Url == OAuthAppUrl {
			token = auth.Token
		}
	}

	if len(token) == 0 {
		authParam := AuthorizationParams{}
		authParam.Scopes = append(authParam.Scopes, "repo")
		authParam.Note = "gh"
		authParam.NoteUrl = OAuthAppUrl

		auth, err := gh.CreateAuthorization(authParam)
		if err != nil {
			return err
		}

		token = auth.Token
	}

	config.Save(config.Config{gh.User, token})

	gh.Authorization = "token " + token

	return nil
}

func (gh *GitHub) performRequest(request *http.Request) (*http.Response, error) {
	if len(gh.Authorization) == 0 {
		err := gh.performBasicAuth()
		if err != nil {
			return nil, err
		}
	}

	request.Header.Set("Authorization", gh.Authorization)

	response, err := gh.httpClient.Do(request)
	if err != nil {
		return response, err
	}

	if response.StatusCode >= 200 && response.StatusCode < 400 {
		return response, err
	}

	err = handleGitHubErrors(response)

	return response, err
}

func handleGitHubErrors(response *http.Response) error {
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		var githubErrors GitHubErrors
		err = json.Unmarshal(body, &githubErrors)
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
			if len(m) > 0 {
				text = text + m + "\n"
			}
		}

		if len(text) == 0 {
			text = githubErrors.Message
		}

		err = errors.New(text)
	}

	return err
}

func unmarshalBody(response *http.Response, v interface{}) error {
	js, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(js, v)
	if err != nil {
		return err
	}

	return nil
}

func (gh *GitHub) httpGet(uri string, extraHeaders map[string]string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", GitHubUrl, uri)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if extraHeaders != nil {
		for h, v := range extraHeaders {
			request.Header.Set(h, v)
		}
	}

	return gh.performRequest(request)
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

	return gh.performRequest(request)
}

type PullRequestParams struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Base  string `json:"base"`
	Head  string `json:"head"`
}

type AuthorizationParams struct {
	Scopes       []string `json:"scopes"`
	Note         string   `json:"note"`
	NoteUrl      string   `json:"note_url"`
	ClientId     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
}

func (gh *GitHub) CreateAuthorization(authParam AuthorizationParams) (*Authorization, error) {
	b, err := json.Marshal(authParam)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(b)
	response, err := gh.httpPost("/authorizations", nil, buffer)

	var auth Authorization
	err = unmarshalBody(response, &auth)
	if err != nil {
		return nil, err
	}

	return &auth, nil
}

type PullRequestResponse struct {
	Url      string `json:"url"`
	HtmlUrl  string `json:"html_url"`
	DiffUrl  string `json:"diff_url"`
	PatchUrl string `json:"patch_url"`
	IssueUrl string `json:"issue_url"`
}

func (gh *GitHub) CreatePullRequest(project *Project, params PullRequestParams) (*PullRequestResponse, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(b)
	url := fmt.Sprintf("/repos/%s/%s/pulls", project.Owner, project.Name)
	response, err := gh.httpPost(url, nil, buffer)
	if err != nil {
		return nil, err
	}

	var pullRequestResponse PullRequestResponse
	err = unmarshalBody(response, &pullRequestResponse)
	if err != nil {
		return nil, err
	}

	return &pullRequestResponse, nil
}
