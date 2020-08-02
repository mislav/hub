package git

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/github/hub/v2/internal/assert"
)

func TestSSHConfigReader_Read(t *testing.T) {
	f, _ := ioutil.TempFile("", "ssh-config")
	c := `Host github.com
  Hostname ssh.github.com
  Port 443

	host other
	Hostname 10.0.0.1
	`

	ioutil.WriteFile(f.Name(), []byte(c), os.ModePerm)

	r := &SSHConfigReader{[]string{f.Name()}}
	sc := r.Read()
	assert.Equal(t, "ssh.github.com", sc["github.com"])
}

func TestSSHConfigReader_ExpandTokens(t *testing.T) {
	f, _ := ioutil.TempFile("", "ssh-config")
	c := `Host github.com example.org
  Hostname 1-%h-2-%%h-3-%h-%%
	`

	ioutil.WriteFile(f.Name(), []byte(c), os.ModePerm)

	r := &SSHConfigReader{[]string{f.Name()}}
	sc := r.Read()
	assert.Equal(t, "1-github.com-2-%h-3-github.com-%", sc["github.com"])
	assert.Equal(t, "1-example.org-2-%h-3-example.org-%", sc["example.org"])
}
