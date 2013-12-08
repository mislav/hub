package github

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestFullBaseAndFullHead(t *testing.T) {
	project := Project{Name: "name", Owner: "owner"}
	repo := Repo{"base", "head", &project}

	assert.Equal(t, "owner:base", repo.FullBase())
	assert.Equal(t, "owner:head", repo.FullHead())
}
