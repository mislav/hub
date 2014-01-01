package commands

import (
	"fmt"
	"github.com/bmizerany/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestUpdater_downloadFile(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/gh.zip", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		fmt.Fprint(w, "1234")
	})

	path, err := downloadFile(fmt.Sprintf("%s/gh.zip", server.URL))
	assert.Equal(t, nil, err)

	content, err := ioutil.ReadFile(path)
	assert.Equal(t, nil, err)
	assert.Equal(t, "1234", string(content))
	assert.Equal(t, "gh.zip", filepath.Base(path))
}
