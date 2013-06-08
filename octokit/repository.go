package octokit

import (
	"fmt"
)

type Repository struct {
	Name     string
	UserName string
}

func (r Repository) String() string {
	return fmt.Sprintf("%s/%s", r.UserName, r.Name)
}
