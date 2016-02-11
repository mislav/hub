package octokit

import (
	"github.com/stretchr/testify/assert"

	"encoding/json"
	"testing"
)

func TestOrganizationService_Get(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/orgs/octokit", "organization", nil)

	organizationResults, result := client.Organization().OrganizationGet(nil, M{"org": "octokit"})

	assert.False(t, result.HasError())
	assert.Equal(t, "octokit", organizationResults.Login)
	assert.Equal(t, 3430433, organizationResults.ID)
}

func TestOrganizationService_Update(t *testing.T) {
	setup()
	defer tearDown()

	input := OrganizationParams{
		BillingEmail: "support@github.com",
		Blog:         "https://github.com/blog",
		Company:      "GitHub",
		Email:        "support@github.com",
		Location:     "San Francisco",
		Name:         "github",
		Description:  "GitHub, the company.",
	}
	wantReqBody, _ := json.Marshal(input)
	stubPatch(t, "/orgs/github", "organization_updated", nil, string(wantReqBody)+"\n", nil)

	organizationResults, result := client.Organization().OrganizationUpdate(nil, input, M{"org": "github"})

	assert.False(t, result.HasError())
	assert.Equal(t, "github", organizationResults.Login)
	assert.Equal(t, 9919, organizationResults.ID)
}

func TestOrganizationService_Repos(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/orgs/rails/repos", "repositories", nil)

	organizationResults, result := client.Organization().OrganizationRepos(nil, M{"org": "rails"})

	assert.False(t, result.HasError())
	assert.Equal(t, 30, len(organizationResults))
	assert.Equal(t, 8514, organizationResults[0].ID)
	assert.Equal(t, "rails", organizationResults[0].Name)
	assert.Equal(t, 13992, organizationResults[1].ID)
	assert.Equal(t, "docrails", organizationResults[1].Name)
}

func TestOrganizationService_Yours(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/user/orgs", "organizations", nil)

	organizationResults, result := client.Organization().YourOrganizations(nil, nil)

	assert.False(t, result.HasError())
	assert.Equal(t, 2, len(organizationResults))
	assert.Equal(t, 1388703, organizationResults[0].ID)
	assert.Equal(t, "acl-services", organizationResults[0].Login)
	assert.Equal(t, 3430433, organizationResults[1].ID)
	assert.Equal(t, "octokit", organizationResults[1].Login)
}

func TestOrganizationService_User(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/users/rails/orgs", "organizations", nil)

	organizationResults, result := client.Organization().UserOrganizations(nil, M{"username": "rails"})

	assert.False(t, result.HasError())
	assert.Equal(t, 2, len(organizationResults))
	assert.Equal(t, 1388703, organizationResults[0].ID)
	assert.Equal(t, "acl-services", organizationResults[0].Login)
	assert.Equal(t, 3430433, organizationResults[1].ID)
	assert.Equal(t, "octokit", organizationResults[1].Login)
}

func TestOrganizationService_Failure(t *testing.T) {
	setup()
	defer tearDown()
	url := Hyperlink("}")
	organizationResultsRepo, result := client.Organization().OrganizationRepos(&url, nil)
	assert.True(t, result.HasError())
	assert.Equal(t, []Repository(nil), organizationResultsRepo)

	organizationResult, result := client.Organization().OrganizationGet(&url, nil)
	assert.True(t, result.HasError())
	assert.Equal(t, Organization{}, organizationResult)

	organizationResult, result = client.Organization().OrganizationUpdate(&url, OrganizationParams{}, nil)
	assert.True(t, result.HasError())
	assert.Equal(t, Organization{}, organizationResult)

	organizationResults, result := client.Organization().YourOrganizations(&url, nil)
	assert.True(t, result.HasError())
	assert.Equal(t, []Organization(nil), organizationResults)

	organizationResults, result = client.Organization().UserOrganizations(&url, nil)
	assert.True(t, result.HasError())
	assert.Equal(t, []Organization(nil), organizationResults)
}
