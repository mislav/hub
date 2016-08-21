package octokit

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/bmizerany/assert"
)

func TestUploadsService_UploadAsset(t *testing.T) {
	setup()
	defer tearDown()

	file, err := ioutil.TempFile("", "octokit-test-upload-")
	assert.Equal(t, nil, err)
	file.WriteString("this is a test")

	fi, err := file.Stat()
	assert.Equal(t, nil, err)
	file.Close()

	mux.HandleFunc("/repos/octokit/Hello-World/releases/123/assets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "text/plain")
		assert.Equal(t, fi.Size(), r.ContentLength)
		respondWithStatus(w, 201)
	})

	link := Hyperlink("/repos/octokit/Hello-World/releases/123/assets{?name}")
	url, err := link.Expand(M{"name": fi.Name()})
	assert.Equal(t, nil, err)

	open, _ := os.Open(file.Name())
	result := client.Uploads(url).UploadAsset(open, "text/plain", fi.Size())
	fmt.Println(result)
	assert.T(t, !result.HasError())

	assert.Equal(t, 201, result.Response.StatusCode)
}
