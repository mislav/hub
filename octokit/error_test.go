package octokit

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestHandleErrors(t *testing.T) {
	body := "{\"message\":\"Invalid request media type (expecting 'text/plain')\"}"
	err := handleErrors([]byte(body))

	assert.Equal(t, "Invalid request media type (expecting 'text/plain')", err.Error())
}
