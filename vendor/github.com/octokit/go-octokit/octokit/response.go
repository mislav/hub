package octokit

import (
	"net/http"

	"github.com/jingweno/go-sawyer"
	"github.com/jingweno/go-sawyer/mediaheader"
	"github.com/jingweno/go-sawyer/mediatype"
)

// Response is a wrapper for a HttpResponse that adds a cleaned form
// of the MeidaType and MediaHeader
type Response struct {
	MediaType   *mediatype.MediaType
	MediaHeader *mediaheader.MediaHeader
	*http.Response
}

// NewResponse unwraps a sawyer Response, producing an error if there
// was one associated in the sawyer response and otherwise creating a
// new Response from the underlying HttpResponse, MediaType and
// MediaHeader
func NewResponse(sawyerResp *sawyer.Response) (resp *Response, err error) {
	if sawyerResp.IsError() {
		err = sawyerResp.ResponseError
		return
	}

	if sawyerResp.IsApiError() {
		err = NewResponseError(sawyerResp)
		return
	}

	resp = &Response{Response: sawyerResp.Response, MediaType: sawyerResp.MediaType, MediaHeader: sawyerResp.MediaHeader}

	return
}
