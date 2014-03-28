package github

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestApiHost_String(t *testing.T) {
	ah := &apiHost{"github.com"}
	assert.Equal(t, "https://api.github.com", ah.String())

	ah = &apiHost{"github.corporate.com"}
	assert.Equal(t, "https://github.corporate.com", ah.String())

	ah = &apiHost{"http://github.corporate.com"}
	assert.Equal(t, "http://github.corporate.com", ah.String())
}
