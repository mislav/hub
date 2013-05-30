package github

import (
	"bytes"
	"encoding/json"
)

type AuthorizationParams struct {
	Scopes       []string `json:"scopes"`
	Note         string   `json:"note"`
	NoteUrl      string   `json:"note_url"`
	ClientId     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
}

type Authorization struct {
	Scopes  []string `json:"scopes"`
	Url     string   `json:"url"`
	App     App      `json:"app"`
	Token   string   `json:"token"`
	Note    string   `json:"note"`
	NoteUrl string   `json:"note_url"`
}

type App struct {
	Url      string `json:"url"`
	Name     string `json:"name"`
	ClientId string `json:"client_id"`
}


func listAuthorizations(gh *GitHub) ([]Authorization, error) {
	response, err := httpGet(gh, "/authorizations", nil)
	if err != nil {
		return nil, err
	}

	var auths []Authorization
	err = unmarshalBody(response, &auths)

	return auths, err
}

func createAuthorization(gh *GitHub, authParam AuthorizationParams) (*Authorization, error) {
	b, err := json.Marshal(authParam)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(b)
	response, err := httpPost(gh, "/authorizations", nil, buffer)

	var auth Authorization
	err = unmarshalBody(response, &auth)
	if err != nil {
		return nil, err
	}

	return &auth, nil
}
