package utils

import (
	"github.com/github/hub/v2/internal/assert"
	"testing"
	"time"
)

func TestSearchBrowserLauncher(t *testing.T) {
	browser := searchBrowserLauncher("darwin")
	assert.Equal(t, "open", browser)

	browser = searchBrowserLauncher("windows")
	assert.Equal(t, "cmd /c start", browser)
}

func TestConcatPaths(t *testing.T) {
	assert.Equal(t, "foo/bar/baz", ConcatPaths("foo", "bar", "baz"))
}

func TestTimeAgo(t *testing.T) {
	timeNow = func() time.Time {
		return time.Date(2018, 10, 28, 14, 34, 58, 651387237, time.UTC)
	}

	now := timeNow()

	secAgo := now.Add(-1 * time.Second)
	actual := TimeAgo(secAgo)
	assert.Equal(t, "now", actual)

	minAgo := now.Add(-1 * time.Minute)
	actual = TimeAgo(minAgo)
	assert.Equal(t, "1 minute ago", actual)

	minsAgo := now.Add(-5 * time.Minute)
	actual = TimeAgo(minsAgo)
	assert.Equal(t, "5 minutes ago", actual)

	hourAgo := now.Add(-1 * time.Hour)
	actual = TimeAgo(hourAgo)
	assert.Equal(t, "1 hour ago", actual)

	hoursAgo := now.Add(-3 * time.Hour)
	actual = TimeAgo(hoursAgo)
	assert.Equal(t, "3 hours ago", actual)

	dayAgo := now.Add(-1 * 24 * time.Hour)
	actual = TimeAgo(dayAgo)
	assert.Equal(t, "1 day ago", actual)

	daysAgo := now.Add(-5 * 24 * time.Hour)
	actual = TimeAgo(daysAgo)
	assert.Equal(t, "5 days ago", actual)

	monthAgo := now.Add(-1 * 24 * 31 * time.Hour)
	actual = TimeAgo(monthAgo)
	assert.Equal(t, "1 month ago", actual)

	monthsAgo := now.Add(-2 * 24 * 31 * time.Hour)
	actual = TimeAgo(monthsAgo)
	assert.Equal(t, "2 months ago", actual)

	yearAgo := now.Add(-1 * 24 * 31 * 12 * time.Hour)
	actual = TimeAgo(yearAgo)
	assert.Equal(t, "1 year ago", actual)

	yearsAgo := now.Add(-2 * 24 * 31 * 12 * time.Hour)
	actual = TimeAgo(yearsAgo)
	assert.Equal(t, "2 years ago", actual)
}
