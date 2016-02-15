package octokit

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailsService_All(t *testing.T) {
	setup()
	defer tearDown()

	link := fmt.Sprintf(`<%s>; rel="next", <%s>; rel="last"`, testURLOf("/user/emails?page=2"), testURLOf("/user/emails?page=3"))
	respHeaderParams := map[string]string{"Link": link}
	stubGet(t, "/user/emails", "emails", respHeaderParams)

	url, _ := EmailUrl.Expand(nil)
	allEmails, result := client.Emails(url).All()

	assert.False(t, result.HasError())
	assert.Len(t, allEmails, 1)

	email := allEmails[0]
	assert.Equal(t, "rz99@cornell.edu", email.Email)
	assert.Equal(t, true, email.Verified)
	assert.Equal(t, true, email.Primary)

	assert.Equal(t, testURLStringOf("/user/emails?page=2"), string(*result.NextPage))
	assert.Equal(t, testURLStringOf("/user/emails?page=3"), string(*result.LastPage))

	nextPageURL, err := result.NextPage.Expand(nil)
	assert.NoError(t, err)

	allEmails, result = client.Emails(nextPageURL).All()
	assert.False(t, result.HasError())
	assert.Len(t, allEmails, 1)
}

func TestEmailsService_Create(t *testing.T) {
	setup()
	defer tearDown()

	url, _ := EmailUrl.Expand(nil)

	params := []string{"test@example.com", "otherTest@example.com"}
	wantReqBody, _ := json.Marshal(params)
	stubPost(t, "/user/emails", "emails", nil, string(wantReqBody)+"\n", nil)

	allEmails, result := client.Emails(url).Create(params)

	assert.False(t, result.HasError())
	assert.Len(t, allEmails, 1)

	email := allEmails[0]
	assert.Equal(t, "rz99@cornell.edu", email.Email)
	assert.Equal(t, true, email.Verified)
	assert.Equal(t, true, email.Primary)
}

func TestEmailsService_Delete(t *testing.T) {
	setup()
	defer tearDown()

	url, _ := EmailUrl.Expand(nil)

	params := []string{"test@example.com", "otherTest@example.com"}
	wantReqBody, _ := json.Marshal(params)
	respHeaderParams := map[string]string{"Content-Type": "application/json"}
	stubDeletewCodewBody(t, "/user/emails", string(wantReqBody)+"\n", respHeaderParams, 204)
	result := client.Emails(url).Delete(params)

	assert.False(t, result.HasError())
	assert.Equal(t, 204, result.Response.StatusCode)
}
