package octokit

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReleasesService_Latest(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/repos/jingweno/gh/releases/latest", "latest_release", nil)

	url, err := ReleasesLatestURL.Expand(M{"owner": "jingweno", "repo": "gh"})
	assert.NoError(t, err)

	release, result := client.Releases(url).Latest()
	assert.False(t, result.HasError())
	assert.Equal(t, 295009, release.ID)
	assert.Equal(t, "v2.1.0", release.TagName)
}

func TestReleasesService_All(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/repos/jingweno/gh/releases", "releases", nil)

	url, err := ReleasesURL.Expand(M{"owner": "jingweno", "repo": "gh"})
	assert.NoError(t, err)

	releases, result := client.Releases(url).All()
	assert.False(t, result.HasError())
	assert.Len(t, releases, 1)

	firstRelease := releases[0]
	assert.Equal(t, 50013, firstRelease.ID)
	assert.Equal(t, "v0.23.0", firstRelease.TagName)
	assert.Equal(t, "master", firstRelease.TargetCommitish)
	assert.Equal(t, "v0.23.0", firstRelease.Name)
	assert.False(t, firstRelease.Draft)
	assert.False(t, firstRelease.Prerelease)
	assert.Equal(t, "* Windows works!: https://github.com/jingweno/gh/commit/6cb80cb09fd9f624a64d85438157955751a9ac70", firstRelease.Body)
	assert.Equal(t, "https://api.github.com/repos/jingweno/gh/releases/50013", firstRelease.URL)
	assert.Equal(t, "https://api.github.com/repos/jingweno/gh/releases/50013/assets", firstRelease.AssetsURL)
	assert.Equal(t, "https://uploads.github.com/repos/jingweno/gh/releases/50013/assets{?name}", string(firstRelease.UploadURL))
	assert.Equal(t, "https://github.com/jingweno/gh/releases/v0.23.0", firstRelease.HTMLURL)
	assert.Equal(t, "2013-09-23 00:59:10 +0000 UTC", firstRelease.CreatedAt.String())
	assert.Equal(t, "2013-09-23 01:07:56 +0000 UTC", firstRelease.PublishedAt.String())

	firstReleaseAssets := firstRelease.Assets
	assert.Len(t, firstReleaseAssets, 8)

	firstAsset := firstReleaseAssets[0]
	assert.Equal(t, 20428, firstAsset.ID)
	assert.Equal(t, "gh_0.23.0-snapshot_amd64.deb", firstAsset.Name)
	assert.Equal(t, "gh_0.23.0-snapshot_amd64.deb", firstAsset.Label)
	assert.Equal(t, "application/x-deb", firstAsset.ContentType)
	assert.Equal(t, "uploaded", firstAsset.State)
	assert.Equal(t, 1562984, firstAsset.Size)
	assert.Equal(t, 0, firstAsset.DownloadCount)
	assert.Equal(t, "https://api.github.com/repos/jingweno/gh/releases/assets/20428", firstAsset.URL)
	assert.Equal(t, "2013-09-23 01:05:20 +0000 UTC", firstAsset.CreatedAt.String())
	assert.Equal(t, "2013-09-23 01:07:56 +0000 UTC", firstAsset.UpdatedAt.String())
}

func TestCreateRelease(t *testing.T) {
	setup()
	defer tearDown()

	url, err := ReleasesURL.Expand(M{"owner": "octokit", "repo": "Hello-World"})
	assert.NoError(t, err)

	params := Release{
		TagName:         "v1.0.0",
		TargetCommitish: "master",
	}
	wantReqBody, _ := json.Marshal(params)
	stubPost(t, "/repos/octokit/Hello-World/releases", "create_release", nil, string(wantReqBody)+"\n", nil)

	release, result := client.Releases(url).Create(params)

	assert.False(t, result.HasError())
	assert.Equal(t, "v1.0.0", release.TagName)
}

func TestUpdateRelease(t *testing.T) {
	setup()
	defer tearDown()

	url, err := ReleasesURL.Expand(M{"owner": "octokit", "repo": "Hello-World", "id": "123"})
	assert.NoError(t, err)

	params := Release{
		TagName:         "v1.0.0",
		TargetCommitish: "master",
	}
	wantReqBody, _ := json.Marshal(params)
	stubPatch(t, "/repos/octokit/Hello-World/releases/123", "create_release", nil, string(wantReqBody)+"\n", nil)

	release, result := client.Releases(url).Update(params)

	assert.False(t, result.HasError())
	assert.Equal(t, "v1.0.0", release.TagName)
}
