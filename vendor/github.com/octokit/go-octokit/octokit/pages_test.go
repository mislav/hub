package octokit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPagesService_PageInfo(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/repos/github/developer.github.com/pages", "pageinfo", nil)

	page, result := client.Pages().PageInfo(&PagesURL, M{"owner": "github",
		"repo": "developer.github.com"})
	assert.False(t, result.HasError())
	assert.Equal(t, "built", page.Status)
	assert.Equal(t, "developer.github.com", page.Cname)
	assert.False(t, page.Custom404)
}

func TestPagesService_PageBuildLatest(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, `/repos/github/developer.github.com/pages/builds/latest`, "page_build", nil)

	build, result := client.Pages().PageBuildLatest(&PagesLatestBuildURL,
		M{"owner": "github", "repo": "developer.github.com"})
	assert.False(t, result.HasError())
	assert.Equal(t, "built", build.Status)
	assert.Equal(t, "351391cdcb88ffae71ec3028c91f375a8036a26b", build.Commit)
	assert.Equal(t, 1, build.Pusher.ID)
	assert.Equal(t, 2104, build.Duration)
}

func TestPagesService_PageBuilds(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/repos/github/developer.github.com/pages/builds", "page_builds", nil)

	builds, result := client.Pages().PageBuilds(&PagesBuildsURL,
		M{"owner": "github", "repo": "developer.github.com"})
	assert.False(t, result.HasError())
	assert.Equal(t, 1, len(builds))
	assert.Equal(t, "built", builds[0].Status)
	assert.Equal(t, "351391cdcb88ffae71ec3028c91f375a8036a26b", builds[0].Commit)
	assert.Equal(t, 1, builds[0].Pusher.ID)
	assert.Equal(t, 2104, builds[0].Duration)
}

func TestPageService_Failure(t *testing.T) {
	setup()
	defer tearDown()
	url := Hyperlink("}")
	pageResult, result := client.Pages().PageInfo(&url, nil)
	assert.True(t, result.HasError())
	assert.Equal(t, (*PageInfo)(nil), pageResult)

	pageBuildResults, result := client.Pages().PageBuilds(&url, nil)
	assert.True(t, result.HasError())
	assert.Equal(t, []PageBuild(nil), pageBuildResults)

	pageBuildResult, result := client.Pages().PageBuildLatest(&url, nil)
	assert.True(t, result.HasError())
	assert.Equal(t, (*PageBuild)(nil), pageBuildResult)
}
