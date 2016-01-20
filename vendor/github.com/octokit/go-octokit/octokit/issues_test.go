package octokit

import (
	"net/http"
	"testing"
	"time"

	"github.com/bmizerany/assert"
)

func TestIssuesService_All(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/octocat/Hello-World/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("issues.json"))
	})

	url, err := RepoIssuesURL.Expand(M{"owner": "octocat", "repo": "Hello-World"})
	assert.Equal(t, nil, err)

	issues, result := client.Issues(url).All()
	assert.T(t, !result.HasError())
	assert.Equal(t, 1, len(issues))

	issue := issues[0]
	validateIssue(t, issue)
}

func TestIssuesService_One(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/octocat/Hello-World/issues/1347", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("issue.json"))
	})

	url, err := RepoIssuesURL.Expand(M{"owner": "octocat", "repo": "Hello-World", "number": 1347})
	assert.Equal(t, nil, err)

	issue, result := client.Issues(url).One()

	assert.T(t, !result.HasError())
	validateIssue(t, *issue)
}

func TestIssuesService_Create(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/octocat/Hello-World/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, "{\"title\":\"title\",\"body\":\"body\"}\n")
		respondWithJSON(w, loadFixture("issue.json"))
	})

	url, err := RepoIssuesURL.Expand(M{"owner": "octocat", "repo": "Hello-World"})
	assert.Equal(t, nil, err)

	params := IssueParams{
		Title: "title",
		Body:  "body",
	}
	issue, result := client.Issues(url).Create(params)

	assert.T(t, !result.HasError())
	validateIssue(t, *issue)
}

func TestIssuesService_Update(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/octocat/Hello-World/issues/1347", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, "{\"title\":\"title\",\"body\":\"body\"}\n")
		respondWithJSON(w, loadFixture("issue.json"))
	})

	url, err := RepoIssuesURL.Expand(M{"owner": "octocat", "repo": "Hello-World", "number": 1347})
	assert.Equal(t, nil, err)

	params := IssueParams{
		Title: "title",
		Body:  "body",
	}
	issue, result := client.Issues(url).Update(params)

	assert.T(t, !result.HasError())
	validateIssue(t, *issue)
}

func validateIssue(t *testing.T, issue Issue) {

	assert.Equal(t, "https://api.github.com/repos/octocat/Hello-World/issues/1347", issue.URL)
	assert.Equal(t, "https://github.com/octocat/Hello-World/issues/1347", issue.HTMLURL)
	assert.Equal(t, 1347, issue.Number)
	assert.Equal(t, "open", issue.State)
	assert.Equal(t, "Found a bug", issue.Title)
	assert.Equal(t, "I'm having a problem with this.", issue.Body)

	assert.Equal(t, "octocat", issue.User.Login)
	assert.Equal(t, 1, issue.User.ID)
	assert.Equal(t, "https://github.com/images/error/octocat_happy.gif", issue.User.AvatarURL)
	assert.Equal(t, "somehexcode", issue.User.GravatarID)
	assert.Equal(t, "https://api.github.com/users/octocat", issue.User.URL)

	assert.Equal(t, 1, len(issue.Labels))
	assert.Equal(t, "https://api.github.com/repos/octocat/Hello-World/labels/bug", issue.Labels[0].URL)
	assert.Equal(t, "bug", issue.Labels[0].Name)

	assert.Equal(t, "octocat", issue.Assignee.Login)
	assert.Equal(t, 1, issue.Assignee.ID)
	assert.Equal(t, "https://github.com/images/error/octocat_happy.gif", issue.Assignee.AvatarURL)
	assert.Equal(t, "somehexcode", issue.Assignee.GravatarID)
	assert.Equal(t, "https://api.github.com/users/octocat", issue.Assignee.URL)

	assert.Equal(t, "https://api.github.com/repos/octocat/Hello-World/milestones/1", issue.Milestone.URL)
	assert.Equal(t, 1, issue.Milestone.Number)
	assert.Equal(t, "open", issue.Milestone.State)
	assert.Equal(t, "v1.0", issue.Milestone.Title)
	assert.Equal(t, "", issue.Milestone.Description)

	assert.Equal(t, "octocat", issue.Milestone.Creator.Login)
	assert.Equal(t, 1, issue.Milestone.Creator.ID)
	assert.Equal(t, "https://github.com/images/error/octocat_happy.gif", issue.Milestone.Creator.AvatarURL)
	assert.Equal(t, "somehexcode", issue.Milestone.Creator.GravatarID)
	assert.Equal(t, "https://api.github.com/users/octocat", issue.Milestone.Creator.URL)

	assert.Equal(t, 4, issue.Milestone.OpenIssues)
	assert.Equal(t, 8, issue.Milestone.ClosedIssues)
	assert.Equal(t, "2011-04-10 20:09:31 +0000 UTC", issue.Milestone.CreatedAt.String())
	assert.Equal(t, (*time.Time)(nil), issue.Milestone.DueOn)

	assert.Equal(t, 0, issue.Comments)
	assert.Equal(t, "https://github.com/octocat/Hello-World/pull/1347", issue.PullRequest.HTMLURL)
	assert.Equal(t, "https://github.com/octocat/Hello-World/pull/1347.diff", issue.PullRequest.DiffURL)
	assert.Equal(t, "https://github.com/octocat/Hello-World/pull/1347.patch", issue.PullRequest.PatchURL)

	assert.Equal(t, (*time.Time)(nil), issue.ClosedAt)
	assert.Equal(t, "2011-04-22 13:33:48 +0000 UTC", issue.CreatedAt.String())
	assert.Equal(t, "2011-04-22 13:33:48 +0000 UTC", issue.UpdatedAt.String())
}
