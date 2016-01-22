package git

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bmizerany/assert"
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
