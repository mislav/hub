package commands

import (
	"github.com/bmizerany/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAssetsDirWithoutFlag(t *testing.T) {
	dir := createTempDir(t)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Chdir(pwd)
		os.RemoveAll(dir)
	}()

	os.Chdir(dir)

	tagDir := filepath.Join(dir, "releases", "v1.0.0")
	assertAssetsDirSelected(t, tagDir, "")
}

func TestAssetsDirWithFlag(t *testing.T) {
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	tagDir := filepath.Join(dir, "releases", "v1.0.0")
	assertAssetsDirSelected(t, tagDir, tagDir)
}

func assertAssetsDirSelected(t *testing.T, expectedDir, flagDir string) {
	assets, err := getAssetsDirectory(flagDir, "v1.0.0")
	assert.NotEqual(t, nil, err) // Error if it doesn't exist

	os.MkdirAll(expectedDir, 0755)
	assets, err = getAssetsDirectory(flagDir, "v1.0.0")
	assert.NotEqual(t, nil, err) // Error if it's empty

	ioutil.TempFile(expectedDir, "gh-test")
	assets, err = getAssetsDirectory(flagDir, "v1.0.0")

	fiExpected, err := os.Stat(expectedDir)
	fiAssets, err := os.Stat(assets)

	assert.Equal(t, nil, err)
	assert.T(t, os.SameFile(fiExpected, fiAssets))
}

func createTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "gh-test-")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}
