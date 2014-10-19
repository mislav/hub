package octokit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/bmizerany/assert"
)

func TestRepositoresService_One(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/jingweno/octokat", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("repository.json"))
	})

	url, err := RepositoryURL.Expand(M{"owner": "jingweno", "repo": "octokat"})
	assert.Equal(t, nil, err)

	repo, result := client.Repositories(url).One()

	assert.T(t, !result.HasError())
	assert.Equal(t, 10575811, repo.ID)
	assert.Equal(t, "octokat", repo.Name)
	assert.Equal(t, "jingweno/octokat", repo.FullName)
	assert.T(t, !repo.Private)
	assert.T(t, !repo.Fork)
	assert.Equal(t, "https://api.github.com/repos/jingweno/octokat", repo.URL)
	assert.Equal(t, "https://github.com/jingweno/octokat", repo.HTMLURL)
	assert.Equal(t, "https://github.com/jingweno/octokat.git", repo.CloneURL)
	assert.Equal(t, "git://github.com/jingweno/octokat.git", repo.GitURL)
	assert.Equal(t, "git@github.com:jingweno/octokat.git", repo.SSHURL)
	assert.Equal(t, "master", repo.MasterBranch)
	assert.T(t, !repo.Permissions.Admin)
	assert.T(t, !repo.Permissions.Push)
	assert.T(t, repo.Permissions.Pull)
}

func TestRepositoresService_All(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/orgs/rails/repos", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		header := w.Header()
		link := fmt.Sprintf(`<%s>; rel="next", <%s>; rel="last"`, testURLOf("organizations/4223/repos?page=2"), testURLOf("organizations/4223/repos?page=3"))
		header.Set("Link", link)

		respondWithJSON(w, loadFixture("repositories.json"))
	})

	url, err := OrgRepositoriesURL.Expand(M{"org": "rails"})
	assert.Equal(t, nil, err)

	repos, result := client.Repositories(url).All()

	assert.T(t, !result.HasError())
	assert.Equal(t, 30, len(repos))
	assert.Equal(t, testURLStringOf("organizations/4223/repos?page=2"), string(*result.NextPage))
	assert.Equal(t, testURLStringOf("organizations/4223/repos?page=3"), string(*result.LastPage))
}

func TestRepositoresService_Create(t *testing.T) {
	setup()
	defer tearDown()

	params := Repository{}
	params.Name = "Hello-World"
	params.Description = "This is your first repo"
	params.Homepage = "https://github.com"
	params.Private = false
	params.HasIssues = true
	params.HasWiki = true
	params.HasDownloads = true

	mux.HandleFunc("/user/repos", func(w http.ResponseWriter, r *http.Request) {
		var repoParams Repository
		json.NewDecoder(r.Body).Decode(&repoParams)
		assert.Equal(t, params.Name, repoParams.Name)
		assert.Equal(t, params.Description, repoParams.Description)
		assert.Equal(t, params.Homepage, repoParams.Homepage)
		assert.Equal(t, params.Private, repoParams.Private)
		assert.Equal(t, params.HasIssues, repoParams.HasIssues)
		assert.Equal(t, params.HasWiki, repoParams.HasWiki)
		assert.Equal(t, params.HasDownloads, repoParams.HasDownloads)

		testMethod(t, r, "POST")
		respondWithJSON(w, loadFixture("create_repository.json"))
	})

	url, err := UserRepositoriesURL.Expand(nil)
	assert.Equal(t, nil, err)

	repo, result := client.Repositories(url).Create(params)

	assert.T(t, !result.HasError())
	assert.Equal(t, 1296269, repo.ID)
	assert.Equal(t, "Hello-World", repo.Name)
	assert.Equal(t, "octocat/Hello-World", repo.FullName)
	assert.Equal(t, "This is your first repo", repo.Description)
	assert.T(t, !repo.Private)
	assert.T(t, repo.Fork)
	assert.Equal(t, "https://api.github.com/repos/octocat/Hello-World", repo.URL)
	assert.Equal(t, "https://github.com/octocat/Hello-World", repo.HTMLURL)
	assert.Equal(t, "https://github.com/octocat/Hello-World.git", repo.CloneURL)
	assert.Equal(t, "git://github.com/octocat/Hello-World.git", repo.GitURL)
	assert.Equal(t, "git@github.com:octocat/Hello-World.git", repo.SSHURL)
	assert.Equal(t, "master", repo.MasterBranch)
}

func TestRepositoresService_CreateFork(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/jingweno/octokat/forks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, "{\"organization\":\"github\"}\n")
		respondWithJSON(w, loadFixture("create_repository.json"))
	})

	url, err := ForksURL.Expand(M{"owner": "jingweno", "repo": "octokat"})
	assert.Equal(t, nil, err)

	repo, result := client.Repositories(url).Create(M{"organization": "github"})

	assert.T(t, !result.HasError())
	assert.Equal(t, 1296269, repo.ID)
	assert.Equal(t, "Hello-World", repo.Name)
	assert.Equal(t, "octocat/Hello-World", repo.FullName)
	assert.Equal(t, "This is your first repo", repo.Description)
	assert.T(t, !repo.Private)
	assert.T(t, repo.Fork)
	assert.Equal(t, "https://api.github.com/repos/octocat/Hello-World", repo.URL)
	assert.Equal(t, "https://github.com/octocat/Hello-World", repo.HTMLURL)
	assert.Equal(t, "https://github.com/octocat/Hello-World.git", repo.CloneURL)
	assert.Equal(t, "git://github.com/octocat/Hello-World.git", repo.GitURL)
	assert.Equal(t, "git@github.com:octocat/Hello-World.git", repo.SSHURL)
	assert.Equal(t, "master", repo.MasterBranch)
}
