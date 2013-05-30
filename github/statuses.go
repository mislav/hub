package github

import (
	"fmt"
	"time"
)

type Status struct {
	CreatedAt   time.Time `josn:"created_at"`
	State       string    `json:"state"`
	TargetUrl   string    `json:"target_url"`
	Description string    `json:"description"`
}

func listStatuses(gh *GitHub, ref string) ([]Status, error) {
	project := gh.project
	url := fmt.Sprintf("/repos/%s/%s/statuses/%s", project.Owner, project.Name, ref)
	response, err := httpGet(gh, url, nil)
	if err != nil {
		return nil, err
	}

	var statuses []Status
	err = unmarshalBody(response, &statuses)
	if err != nil {
		return nil, err
	}

	return statuses, nil
}
