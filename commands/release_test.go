package commands

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/github/hub/fixtures"
)

func TestAssetFinder_Find(t *testing.T) {
	finder := assetFinder{}

	paths, err := finder.Find(fixtures.Path("release_dir", "file1"))
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(paths))

	paths, err = finder.Find(fixtures.Path("release_dir", "dir"))
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(paths))
}

func TestAssetUploader_detectContentType(t *testing.T) {
	u := &assetUploader{}
	ct, err := u.detectContentType(fixtures.Path("release_dir", "file1"))

	assert.Equal(t, nil, err)
	assert.Equal(t, "text/plain", ct)
}
