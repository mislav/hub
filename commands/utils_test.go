package commands

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestRemoveItem(t *testing.T) {
	slice := []string{"1", "2", "3"}
	slice, item := removeItem(slice, 1)

	assert.Equal(t, "2", item)
	assert.Equal(t, 2, len(slice))
	assert.Equal(t, "1", slice[0])
	assert.Equal(t, "3", slice[1])
}
