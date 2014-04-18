package fixtures

import (
	"io/ioutil"
	"os"
)

type TestConfigs struct {
	Path string
}

func (c *TestConfigs) TearDown() {
	os.Setenv("GH_CONFIG", "")
	os.RemoveAll(c.Path)
}

func SetupTestConfigs() *TestConfigs {
	file, _ := ioutil.TempFile("", "test-gh-config-")

	content := `[[hosts]]
  host = "github.com"
  user = "jingweno"
  access_token = "123"
  protocol = "http"`
	ioutil.WriteFile(file.Name(), []byte(content), os.ModePerm)
	os.Setenv("GH_CONFIG", file.Name())

	return &TestConfigs{file.Name()}
}
