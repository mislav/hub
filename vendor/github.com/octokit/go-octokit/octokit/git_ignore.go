package octokit

// GitIgnoreURL is an address for accessing various templates to apply
// to a repository upon creation.
//
// https://developer.github.com/v3/gitignore/
var GitIgnoreURL = Hyperlink("/gitignore/templates{/name}")

// GitIgnore creates a GitIgnoreService to access gitignore templates
//
// https://developer.github.com/v3/gitignore/
func (c *Client) GitIgnore() *GitIgnoreService {
	return &GitIgnoreService{client: c}
}

// A service to return gitignore templates
type GitIgnoreService struct {
	client *Client
}

// All gets a list all the available templates
//
// https://developer.github.com/v3/gitignore/#listing-available-templates
func (s *GitIgnoreService) All(uri *Hyperlink) (templates []string, result *Result) {
	url, err := ExpandWithDefault(uri, &GitIgnoreURL, nil)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = s.client.get(url, &templates)
	return
}

// One gets a specific gitignore template based on the passed url
//
// https://developer.github.com/v3/gitignore/#get-a-single-template
func (s *GitIgnoreService) One(uri *Hyperlink, uriParams M) (template *GitIgnoreTemplate, result *Result) {
	url, err := ExpandWithDefault(uri, &GitIgnoreURL, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = s.client.get(url, &template)
	return
}

//GitIgnoreTemplate is a representation of a given template returned by the service
type GitIgnoreTemplate struct {
	Name   string `json:"name,omitempty"`
	Source string `json:"source,omitempty"`
}
