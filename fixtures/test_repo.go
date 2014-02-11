package fixtures

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/github/hub/cmd"
)

var pwd string

func init() {
	// need to cache pwd before all tests run
	// pwd is changed to the bin dir in the tmp folder during test run
	pwd, _ = os.Getwd()
}

type TestRepo struct {
	pwd    string
	dir    string
	Remote string
}

func (r *TestRepo) Setup() (err error) {
	r.dir, err = ioutil.TempDir("", "test-repo")
	if err != nil {
		return
	}
	targetPath := filepath.Join(r.dir, "test.git")

	err = r.clone(r.Remote, targetPath)
	if err != nil {
		return
	}

	return os.Chdir(targetPath)
}

func (r *TestRepo) clone(repo, dir string) error {
	cmd := cmd.New("git").WithArgs("clone", repo, dir)
	output, err := cmd.ExecOutput()
	if err != nil {
		err = fmt.Errorf("error cloning %s to %s: %s", repo, dir, output)
	}

	return err
}

func (r *TestRepo) TearDown() error {
	err := os.Remove(r.dir)
	if err != nil {
		return err
	}

	err = os.Chdir(r.pwd)
	if err != nil {
		return err
	}

	return nil
}

func SetupTestRepo() *TestRepo {
	remotePath := filepath.Join(pwd, "..", "fixtures", "test.git")
	repo := &TestRepo{pwd: pwd, Remote: remotePath}
	err := repo.Setup()
	if err != nil {
		panic(err)
	}

	return repo
}
