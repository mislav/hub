package fixtures

import (
	"os"
	"path/filepath"
)

type TestRepo struct {
	Pwd string
}

func (g *TestRepo) Setup() {
	g.Pwd, _ = os.Getwd()
	fixturePath := filepath.Join(g.Pwd, "..", "fixtures", "test.git")
	err := os.Chdir(fixturePath)
	if err != nil {
		panic(err)
	}
}

func (g *TestRepo) TearDown() {
	err := os.Chdir(g.Pwd)
	if err != nil {
		panic(err)
	}
}

func SetupTestRepo() *TestRepo {
	repo := &TestRepo{}
	repo.Setup()

	return repo
}
