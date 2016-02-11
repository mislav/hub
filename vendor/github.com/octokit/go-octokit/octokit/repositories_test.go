package octokit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepositoresService_One(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/repos/jingweno/octokat", "repository", nil)

	repo, result := client.Repositories().One(nil, M{"owner": "jingweno",
		"repo": "octokat"})

	assert.False(t, result.HasError())
	assert.Equal(t, 10575811, repo.ID)
	assert.Equal(t, "octokat", repo.Name)
	assert.EqualValues(t, "jingweno/octokat", repo.FullName)
	assert.False(t, repo.Private)
	assert.False(t, repo.Fork)
	assert.EqualValues(t, "https://api.github.com/repos/jingweno/octokat",
		repo.URL)
	assert.EqualValues(t, "https://github.com/jingweno/octokat",
		repo.HTMLURL)
	assert.EqualValues(t, "https://github.com/jingweno/octokat.git",
		repo.CloneURL)
	assert.EqualValues(t, "git://github.com/jingweno/octokat.git",
		repo.GitURL)
	assert.EqualValues(t, "git@github.com:jingweno/octokat.git",
		repo.SSHURL)
	assert.Equal(t, 79, repo.StargazersCount)
	assert.EqualValues(t, "master", repo.MasterBranch)
	assert.False(t, repo.Permissions.Admin)
	assert.False(t, repo.Permissions.Push)
	assert.True(t, repo.Permissions.Pull)
}

func TestRepositoresService_All(t *testing.T) {
	setup()
	defer tearDown()

	link := fmt.Sprintf(`<%s>; rel="next", <%s>; rel="last"`,
		testURLOf("organizations/4223/repos?page=2"),
		testURLOf("organizations/4223/repos?page=3"))
	stubGet(t, "/orgs/rails/repos", "repositories",
		map[string]string{"Link": link})

	repos, result := client.Repositories().All(&OrgRepositoriesURL,
		M{"org": "rails"})

	fmt.Println(result.Error())
	assert.False(t, result.HasError())
	assert.Len(t, repos, 30)
	assert.Equal(t, testURLStringOf("organizations/4223/repos?page=2"),
		string(*result.NextPage))
	assert.Equal(t, testURLStringOf("organizations/4223/repos?page=3"),
		string(*result.LastPage))
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

	mux.HandleFunc("/user/repos", func(w http.ResponseWriter,
		r *http.Request) {
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

	repo, result := client.Repositories().Create(&UserRepositoriesURL, nil,
		params)

	assert.False(t, result.HasError())
	assert.Equal(t, 1296269, repo.ID)
	assert.EqualValues(t, "Hello-World", repo.Name)
	assert.EqualValues(t, "octocat/Hello-World", repo.FullName)
	assert.EqualValues(t, "This is your first repo", repo.Description)
	assert.False(t, repo.Private)
	assert.True(t, repo.Fork)
	assert.EqualValues(t,
		"https://api.github.com/repos/octocat/Hello-World", repo.URL)
	assert.EqualValues(t, "https://github.com/octocat/Hello-World",
		repo.HTMLURL)
	assert.EqualValues(t, "https://github.com/octocat/Hello-World.git",
		repo.CloneURL)
	assert.EqualValues(t, "git://github.com/octocat/Hello-World.git",
		repo.GitURL)
	assert.EqualValues(t, "git@github.com:octocat/Hello-World.git",
		repo.SSHURL)
	assert.EqualValues(t, "master", repo.MasterBranch)
}

func TestRepositoresService_CreateFork(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/jingweno/octokat/forks", func(
		w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, "{\"organization\":\"github\"}\n")
		respondWithJSON(w, loadFixture("create_repository.json"))
	})

	params := M{"owner": "jingweno", "repo": "octokat"}

	repo, result := client.Repositories().Create(&ForksURL, params,
		M{"organization": "github"})

	assert.False(t, result.HasError())
	assert.Equal(t, 1296269, repo.ID)
	assert.EqualValues(t, "Hello-World", repo.Name)
	assert.EqualValues(t, "octocat/Hello-World", repo.FullName)
	assert.EqualValues(t, "This is your first repo", repo.Description)
	assert.False(t, repo.Private)
	assert.True(t, repo.Fork)
	assert.EqualValues(t,
		"https://api.github.com/repos/octocat/Hello-World", repo.URL)
	assert.EqualValues(t, "https://github.com/octocat/Hello-World",
		repo.HTMLURL)
	assert.EqualValues(t, "https://github.com/octocat/Hello-World.git",
		repo.CloneURL)
	assert.EqualValues(t, "git://github.com/octocat/Hello-World.git",
		repo.GitURL)
	assert.EqualValues(t, "git@github.com:octocat/Hello-World.git",
		repo.SSHURL)
	assert.EqualValues(t, "master", repo.MasterBranch)
}
