//+build !windows
package github

import (
	"fmt"
	"github.com/bmizerany/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestNewHttpClient_UnixSocket(t *testing.T) {
	sock := "/tmp/hub-go.sock"
	s := setupTestServer(sock)
	defer s.Close()

	s.HandleFunc("/unix-socket", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("unix-socket-works"))
	})
	c := newHttpClient("", false, sock)
	resp, err := c.Get(fmt.Sprintf("%s/unix-socket", s.URL.String()))
	assert.Equal(t, nil, err)
	result, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "unix-socket-works", string(result))
}
