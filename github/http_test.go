package github

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestUnmarshalBody(t *testing.T) {
	body := `
[
	{
    "scopes": [
      "repo"
    ],
    "id": 2676297,
    "url": "https://api.github.com/authorizations/2676297",
    "app": {
      "client_id": "9a78b46ea6940243193d",
      "url": "http://owenou.com/gh1",
      "name": "gh1 (API)"
    },
    "token": "4e4428e025a6835f5350f7e97eac1af7a5d11fe2",
    "note": "gh1",
    "note_url": "http://owenou.com/gh1",
    "created_at": "2013-05-29T23:06:54Z",
    "updated_at": "2013-05-29T23:06:54Z"
  },
  {
    "scopes": [
      "repo"
    ],
    "id": 2676299,
    "url": "https://api.github.com/authorizations/2676299",
    "app": {
      "client_id": "da64c1e3003dff224bb0",
      "url": "http://owenou.com/gh2",
      "name": "gh2 (API)"
    },
    "token": "c8310a75f12db067af3280b4c0db3e6f019dacae",
    "note": "gh2",
    "note_url": "http://owenou.com/gh2",
    "created_at": "2013-05-29T23:07:19Z",
    "updated_at": "2013-05-29T23:07:19Z"
  }
]
`
	var auths []Authorization
	err := unmarshal([]byte(body), &auths)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(auths))
	assert.Equal(t, "gh1", auths[0].Note)
	assert.Equal(t, "http://owenou.com/gh1", auths[0].NoteUrl)
}
