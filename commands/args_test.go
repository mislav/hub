package commands

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestRemove(t *testing.T) {
	args := NewArgs([]string{"1", "2", "3"})
	item := args.Remove(1)

	assert.Equal(t, "2", item)
	assert.Equal(t, 2, args.Size())
	assert.Equal(t, "1", args.First())
	assert.Equal(t, "3", args.Get(1))
}
