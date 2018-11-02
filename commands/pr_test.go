package commands

import (
	"testing"
)

func TestPR_ListOpenPR(t *testing.T) {
	cmds := Command{}
	args := Args{}

	listPulls(&cmds, &args)
}
