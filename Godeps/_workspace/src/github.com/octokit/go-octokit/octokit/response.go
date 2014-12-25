package octokit

import (
	"net/http"

	"github.com/jingweno/go-sawyer"
	"github.com/jingweno/go-sawyer/mediaheader"
	"github.com/jingweno/go-sawyer/mediatype"
)

type Response struct {
	MediaType   *mediatype.MediaType
	MediaHeader *mediaheader.MediaHeader
	*http.Response
}

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
