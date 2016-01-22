package mediaheader

import (
	"github.com/bmizerany/assert"
	"net/http"
	"testing"
)

func TestDecoder_Decode(t *testing.T) {
	link := `<https://api.github.com/user/repos?page=3&per_page=100>; rel="next", <https://api.github.com/user/repos?page=50&per_page=100>; rel="last"`
	header := http.Header{}
	header.Add("Link", link)
	decoder := Decoder{}
	mediaHeader := decoder.Decode(header)

	assert.Equal(t, "https://api.github.com/user/repos?page=3&per_page=100", string(mediaHeader.Relations["next"]))
	assert.Equal(t, "https://api.github.com/user/repos?page=50&per_page=100", string(mediaHeader.Relations["last"]))
}
