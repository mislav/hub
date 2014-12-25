package sawyer

import (
	"errors"
	"net/http"

	"github.com/jingweno/go-sawyer/mediaheader"
	"github.com/jingweno/go-sawyer/mediatype"
)

type Response struct {
	ResponseError error
	MediaType     *mediatype.MediaType
	MediaHeader   *mediaheader.MediaHeader
	isApiError    bool
	BodyClosed    bool
	*http.Response
}

func (r *Response) AnyError() bool {
	return r.IsError() || r.IsApiError()
}

func (r *Response) IsError() bool {
	return r.ResponseError != nil
}

func (r *Response) IsApiError() bool {
	return r.isApiError
}

func (r *Response) Error() string {
	if r.ResponseError != nil {
		return r.ResponseError.Error()
	}
	return ""
}

func (r *Response) Decode(resource interface{}) error {
	if r.MediaType == nil {
		return errors.New("No media type for this response")
	}

	if resource == nil || r.ResponseError != nil || r.BodyClosed {
		return r.ResponseError
	}

	defer r.Body.Close()
	r.BodyClosed = true

	dec, err := r.MediaType.Decoder(r.Body)
	if err != nil {
		r.ResponseError = err
	} else {
		r.ResponseError = dec.Decode(resource)
	}
	return r.ResponseError
}

func (r *Response) decode(output interface{}) {
	if !r.isApiError {
		r.Decode(output)
	}
}

func ResponseError(err error) *Response {
	return &Response{ResponseError: err, BodyClosed: true}
}

func UseApiError(status int) bool {
	switch {
	case status > 199 && status < 300:
		return false
	case status == 304:
		return false
	case status == 0:
		return false
	}
	return true
}

func mediaType(res *http.Response) (*mediatype.MediaType, error) {
	if ctype := res.Header.Get(ctypeHeader); len(ctype) > 0 {
		return mediatype.Parse(ctype)
	}
	return nil, nil
}
