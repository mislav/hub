package fixtures

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/github/hub/v2/cmd"
)

type TestRepo struct {
	Remote   string
	TearDown func()
}

func (r *TestRepo) AddRemote(name, url, pushURL string) {
	add := cmd.New("git").WithArgs("remote", "add", name, url)
	if _, err := add.CombinedOutput(); err != nil {
		panic(err)
	}
	if pushURL != "" {
		set := cmd.New("git").WithArgs("remote", "set-url", "--push", name, pushURL)
		if _, err := set.CombinedOutput(); err != nil {
			panic(err)
		}
	}
}

func (r *TestRepo) AddFile(filePath string, content string) {
	path := filepath.Join(os.Getenv("HOME"), filePath)
	err := os.MkdirAll(filepath.Dir(path), 0771)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(path, []byte(content), os.ModePerm)
}

func SetupTestRepo() *TestRepo {
	pwd, _ := os.Getwd()
	oldEnv := make(map[string]string)
	overrideEnv := func(name, value string) {
		oldEnv[name] = os.Getenv(name)
		os.Setenv(name, value)
	}

	remotePath := filepath.Join(pwd, "..", "fixtures", "test.git")
	home, err := ioutil.TempDir("", "test-repo")
	if err != nil {
		panic(err)
	}

	overrideEnv("HOME", home)
	overrideEnv("XDG_CONFIG_HOME", "")
	overrideEnv("XDG_CONFIG_DIRS", "")

	targetPath := filepath.Join(home, "test.git")
	cmd := cmd.New("git").WithArgs("clone", remotePath, targetPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Errorf("error running %s\n%s\n%s", cmd, err, output))
	}

	if err = os.Chdir(targetPath); err != nil {
		panic(err)
	}

	tearDown := func() {
		if err := os.Chdir(pwd); err != nil {
			panic(err)
		}
		for name, value := range oldEnv {
			os.Setenv(name, value)
		}
		if err = os.RemoveAll(home); err != nil {
			panic(err)
		}
	}

	return &TestRepo{Remote: remotePath, TearDown: tearDown}
}
