package octokit

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPublicKeysService_AllKeys(t *testing.T) {
	setup()
	defer tearDown()

	link := fmt.Sprintf(`<%s>; rel="next", <%s>; rel="last"`, testURLOf("/users/obsc/keys?page=2"), testURLOf("/users/obsc/keys?page=3"))
	stubGet(t, "/users/obsc/keys", "public_keys", map[string]string{"Link": link})

	keys, result := client.PublicKeys().All(&PublicKeyUrl, M{"user": "obsc"})
	assert.False(t, result.HasError())
	assert.Len(t, keys, 1)

	key := keys[0]
	assert.Equal(t, 8675080, key.Id)
	assert.Equal(t, "ssh-rsa AAA...", key.Key)

	assert.Equal(t, testURLStringOf("/users/obsc/keys?page=2"), string(*result.NextPage))
	assert.Equal(t, testURLStringOf("/users/obsc/keys?page=3"), string(*result.LastPage))

	validateNextPage_PublicKeys(t, result)
}

func TestPublicKeysService_AllKeysCurrent(t *testing.T) {
	setup()
	defer tearDown()

	link := fmt.Sprintf(`<%s>; rel="next", <%s>; rel="last"`, testURLOf("/user/keys?page=2"), testURLOf("/user/keys?page=3"))
	stubGet(t, "/user/keys", "keys", map[string]string{"Link": link})

	keys, result := client.PublicKeys().All(nil, nil)
	assert.False(t, result.HasError())
	assert.Len(t, keys, 1)

	validateKey(t, keys[0])

	assert.Equal(t, testURLStringOf("/user/keys?page=2"), string(*result.NextPage))
	assert.Equal(t, testURLStringOf("/user/keys?page=3"), string(*result.LastPage))

	validateNextPage_PublicKeys(t, result)
}

func TestPublicKeysService_OneKey(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/user/keys/8675080", "key", nil)

	key, result := client.PublicKeys().One(nil, M{"id": 8675080})
	assert.False(t, result.HasError())

	validateKey(t, *key)
}

func TestPublicKeysService_Create(t *testing.T) {
	setup()
	defer tearDown()

	params := Key{Title: "aKey", Key: "ssh-rsa AAA..."}
	wantReqBody, _ := json.Marshal(params)
	stubPost(t, "/user/keys", "key", nil, string(wantReqBody)+"\n", nil)

	key, result := client.PublicKeys().Create(nil, nil, params)
	assert.False(t, result.HasError())

	validateKey(t, *key)
}

func TestPublicKeysService_Delete(t *testing.T) {
	setup()
	defer tearDown()

	respHeaderParams := map[string]string{"Content-Type": "application/json"}
	stubDeletewCode(t, "/user/keys/8675080", respHeaderParams, 204)

	success, result := client.PublicKeys().Delete(nil, M{"id": 8675080})
	assert.False(t, result.HasError())

	assert.True(t, success)
}

func TestPublicKeysService_Failure(t *testing.T) {
	setup()
	defer tearDown()

	url := Hyperlink("}")
	keys, result := client.PublicKeys().All(&url, nil)
	assert.True(t, result.HasError())
	assert.Len(t, keys, 0)

	key, result := client.PublicKeys().One(&url, nil)
	assert.True(t, result.HasError())
	assert.Nil(t, key)

	key, result = client.PublicKeys().Create(&url, nil, nil)
	assert.True(t, result.HasError())
	assert.Nil(t, key)

	success, result := client.PublicKeys().Delete(&url, nil)
	assert.True(t, result.HasError())
	assert.False(t, success)
}

func validateKey(t *testing.T, key Key) {
	testTime, _ := time.Parse("2006-01-02T15:04:05Z", "2014-07-23T08:42:44Z")

	assert.Equal(t, 8675080, key.Id)
	assert.Equal(t, "ssh-rsa AAA...", key.Key)
	assert.Equal(t, "https://api.github.com/user/keys/8675080", key.URL)
	assert.Equal(t, "aKey", key.Title)
	assert.Equal(t, true, key.Verified)
	assert.Equal(t, &testTime, key.CreatedAt)
}

func validateNextPage_PublicKeys(t *testing.T, result *Result) {
	keys, result := client.PublicKeys().All(result.NextPage, nil)
	assert.False(t, result.HasError())
	assert.Len(t, keys, 1)
}
