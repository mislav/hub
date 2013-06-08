package octokit

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

type AuthorizationParams struct {
	Scopes       []string `json:"scopes"`
	Note         string   `json:"note"`
	NoteUrl      string   `json:"note_url"`
	ClientId     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
}

func (c *Client) Authorizations() ([]Authorization, error) {
	body, err := c.get("authorizations", nil)
	if err != nil {
		return nil, err
	}

	var auths []Authorization
	err = jsonUnmarshal(body, &auths)
	if err != nil {
		return nil, err
	}

	return auths, nil
}

func (c *Client) CreatedAuthorization(params AuthorizationParams) (*Authorization, error) {
	body, err := c.postWithParams("authorizations", nil, params)
	if err != nil {
		return nil, err
	}

	var auth Authorization
	err = jsonUnmarshal(body, &auth)
	if err != nil {
		return nil, err
	}

	return &auth, nil
}
