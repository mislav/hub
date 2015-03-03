package fixtures

import (
	"io/ioutil"
	"os"
)

type TestConfigs struct {
	Path string
}

func (c *TestConfigs) TearDown() {
	os.Setenv("HUB_CONFIG", "")
	os.RemoveAll(c.Path)
}

func SetupTomlTestConfig() *TestConfigs {
	file, _ := ioutil.TempFile("", "test-gh-config-")

	content := `[[hosts]]
  host = "github.com"
  user = "jingweno"
  access_token = "123"
  protocol = "http"`
	ioutil.WriteFile(file.Name(), []byte(content), os.ModePerm)
	os.Setenv("HUB_CONFIG", file.Name())

	return &TestConfigs{file.Name()}
}

func SetupTestConfigs() *TestConfigs {
	file, _ := ioutil.TempFile("", "test-gh-config-")

	content := `---
github.com:
- user: jingweno
  oauth_token: 123
  protocol: http`
	ioutil.WriteFile(file.Name(), []byte(content), os.ModePerm)
	os.Setenv("HUB_CONFIG", file.Name())

	return &TestConfigs{file.Name()}
}
