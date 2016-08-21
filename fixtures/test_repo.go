package fixtures

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/github/hub/cmd"
)

var pwd, home string

func init() {
	// caching `pwd` and $HOME and reset them after test repo is teared down
	// `pwd` is changed to the bin temp dir during test run
	pwd, _ = os.Getwd()
	home = os.Getenv("HOME")
}

type TestRepo struct {
	pwd    string
	dir    string
	home   string
	Remote string
}

func (r *TestRepo) Setup() {
	dir, err := ioutil.TempDir("", "test-repo")
	if err != nil {
		panic(err)
	}
	r.dir = dir

	os.Setenv("HOME", r.dir)

	targetPath := filepath.Join(r.dir, "test.git")
	err = r.clone(r.Remote, targetPath)
	if err != nil {
		panic(err)
	}

	err = os.Chdir(targetPath)
	if err != nil {
		panic(err)
	}
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
	path := filepath.Join(r.dir, filePath)
	err := os.MkdirAll(filepath.Dir(path), 0771)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(path, []byte(content), os.ModePerm)
}

func (r *TestRepo) clone(repo, dir string) error {
	cmd := cmd.New("git").WithArgs("clone", repo, dir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("error cloning %s to %s: %s", repo, dir, output)
	}

	return err
}

func (r *TestRepo) TearDown() {
	err := os.Chdir(r.pwd)
	if err != nil {
		panic(err)
	}

	os.Setenv("HOME", r.home)

	err = os.RemoveAll(r.dir)
	if err != nil {
		panic(err)
	}

}

func SetupTestRepo() *TestRepo {
	remotePath := filepath.Join(pwd, "..", "fixtures", "test.git")
	repo := &TestRepo{pwd: pwd, home: home, Remote: remotePath}
	repo.Setup()

	return repo
}
