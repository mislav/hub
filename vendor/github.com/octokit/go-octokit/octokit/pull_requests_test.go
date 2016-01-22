package octokit

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/bmizerany/assert"
)

func TestPullRequestService_One(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/octokit/go-octokit/pulls/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("pull_request.json"))
	})

	url, err := PullRequestsURL.Expand(M{"owner": "octokit", "repo": "go-octokit", "number": 1})
	assert.Equal(t, nil, err)

	pr, result := client.PullRequests(url).One()

	assert.T(t, !result.HasError())
	assert.Equal(t, 1, pr.ChangedFiles)
	assert.Equal(t, 1, pr.Deletions)
	assert.Equal(t, 1, pr.Additions)
	assert.Equal(t, "aafce5dfc44270f35270b24677abbb859b20addf", pr.MergeCommitSha)
	assert.Equal(t, "2013-06-09 00:53:38 +0000 UTC", pr.MergedAt.String())
	assert.Equal(t, "2013-06-09 00:53:38 +0000 UTC", pr.ClosedAt.String())
	assert.Equal(t, "2013-06-19 00:35:24 +0000 UTC", pr.UpdatedAt.String())
	assert.Equal(t, "2013-06-09 00:52:12 +0000 UTC", pr.CreatedAt.String())
	assert.Equal(t, "typo", pr.Body)
	assert.Equal(t, "Update README.md", pr.Title)
	assert.Equal(t, "https://api.github.com/repos/jingweno/octokat/pulls/1", pr.URL)
	assert.Equal(t, 6206442, pr.ID)
	assert.Equal(t, "https://github.com/jingweno/octokat/pull/1", pr.HTMLURL)
	assert.Equal(t, "https://github.com/jingweno/octokat/pull/1.diff", pr.DiffURL)
	assert.Equal(t, "https://github.com/jingweno/octokat/pull/1.patch", pr.PatchURL)
	assert.Equal(t, "https://github.com/jingweno/octokat/pull/1", pr.IssueURL)
	assert.Equal(t, 1, pr.Number)
	assert.Equal(t, "closed", pr.State)
	assert.T(t, nil == pr.Assignee)
	assert.Equal(t, "https://github.com/jingweno/octokat/pull/1/commits", pr.CommitsURL)
	assert.Equal(t, "https://github.com/jingweno/octokat/pull/1/comments", pr.ReviewCommentsURL)
	assert.Equal(t, "/repos/jingweno/octokat/pulls/comments/{number}", pr.ReviewCommentURL)
	assert.Equal(t, "https://api.github.com/repos/jingweno/octokat/issues/1/comments", pr.CommentsURL)
}

func TestPullRequestService_Post(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/octokit/go-octokit/pulls", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r,
			"{\"base\":\"base\",\"head\":\"head\",\"title\":\"title\",\"body\":\"body\"}\n")
		respondWithJSON(w, loadFixture("pull_request.json"))
	})

	url, err := PullRequestsURL.Expand(M{"owner": "octokit", "repo": "go-octokit"})
	assert.Equal(t, nil, err)

	params := PullRequestParams{
		Base:  "base",
		Head:  "head",
		Title: "title",
		Body:  "body",
	}
	pr, result := client.PullRequests(url).Create(params)

	assert.T(t, !result.HasError())
	assert.Equal(t, 1, pr.ChangedFiles)
	assert.Equal(t, 1, pr.Deletions)
	assert.Equal(t, 1, pr.Additions)
	assert.Equal(t, "aafce5dfc44270f35270b24677abbb859b20addf", pr.MergeCommitSha)
	assert.Equal(t, "2013-06-09 00:53:38 +0000 UTC", pr.MergedAt.String())
	assert.Equal(t, "2013-06-09 00:53:38 +0000 UTC", pr.ClosedAt.String())
	assert.Equal(t, "2013-06-19 00:35:24 +0000 UTC", pr.UpdatedAt.String())
	assert.Equal(t, "2013-06-09 00:52:12 +0000 UTC", pr.CreatedAt.String())
	assert.Equal(t, "typo", pr.Body)
	assert.Equal(t, "Update README.md", pr.Title)
	assert.Equal(t, "https://api.github.com/repos/jingweno/octokat/pulls/1", pr.URL)
	assert.Equal(t, 6206442, pr.ID)
	assert.Equal(t, "https://github.com/jingweno/octokat/pull/1", pr.HTMLURL)
	assert.Equal(t, "https://github.com/jingweno/octokat/pull/1.diff", pr.DiffURL)
	assert.Equal(t, "https://github.com/jingweno/octokat/pull/1.patch", pr.PatchURL)
	assert.Equal(t, "https://github.com/jingweno/octokat/pull/1", pr.IssueURL)
	assert.Equal(t, 1, pr.Number)
	assert.Equal(t, "closed", pr.State)
	assert.T(t, nil == pr.Assignee)
	assert.Equal(t, "https://github.com/jingweno/octokat/pull/1/commits", pr.CommitsURL)
	assert.Equal(t, "https://github.com/jingweno/octokat/pull/1/comments", pr.ReviewCommentsURL)
	assert.Equal(t, "/repos/jingweno/octokat/pulls/comments/{number}", pr.ReviewCommentURL)
	assert.Equal(t, "https://api.github.com/repos/jingweno/octokat/issues/1/comments", pr.CommentsURL)
}

func TestPullRequestService_All(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/rails/rails/pulls", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		header := w.Header()
		link := fmt.Sprintf(`<%s>; rel="next", <%s>; rel="last"`, testURLOf("repositories/8514/pulls?page=2"), testURLOf("repositories/8514/pulls?page=14"))
		header.Set("Link", link)
		respondWithJSON(w, loadFixture("pull_requests.json"))
	})

	url, err := PullRequestsURL.Expand(M{"owner": "rails", "repo": "rails"})
	assert.Equal(t, nil, err)

	prs, result := client.PullRequests(url).All()
	assert.T(t, !result.HasError())
	assert.Equal(t, 30, len(prs))
	assert.Equal(t, testURLStringOf("repositories/8514/pulls?page=2"), string(*result.NextPage))
	assert.Equal(t, testURLStringOf("repositories/8514/pulls?page=14"), string(*result.LastPage))
}

func TestPullRequestService_Diff(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/octokit/go-octokit/pulls/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", diffMediaType)
		respondWith(w, "diff --git")
	})

	url, err := PullRequestsURL.Expand(M{"owner": "octokit", "repo": "go-octokit", "number": 1})
	assert.Equal(t, nil, err)

	diff, result := client.PullRequests(url).Diff()

	assert.T(t, !result.HasError())
	content, err := ioutil.ReadAll(diff)
	assert.Equal(t, nil, err)
	assert.Equal(t, "diff --git", string(content))
}

func TestPullRequestService_Patch(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/octokit/go-octokit/pulls/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", patchMediaType)
		respondWith(w, "patches galore")
	})

	url, err := PullRequestsURL.Expand(M{"owner": "octokit", "repo": "go-octokit", "number": 1})
	assert.Equal(t, nil, err)

	patch, result := client.PullRequests(url).Patch()

	assert.T(t, !result.HasError())
	content, err := ioutil.ReadAll(patch)
	assert.Equal(t, nil, err)
	assert.Equal(t, "patches galore", string(content))
}
