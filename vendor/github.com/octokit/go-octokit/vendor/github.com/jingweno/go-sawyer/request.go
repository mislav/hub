package sawyer

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jingweno/go-sawyer/mediaheader"
	"github.com/jingweno/go-sawyer/mediatype"
)

type Request struct {
	Client    *http.Client
	MediaType *mediatype.MediaType
	Query     url.Values
	*http.Request
}

func (c *Client) NewRequest(rawurl string) (*Request, error) {
	u, err := c.ResolveReferenceString(rawurl)
	if err != nil {
		return nil, err
	}

	httpreq, err := http.NewRequest(GetMethod, u, nil)
	for key, _ := range c.Header {
		httpreq.Header.Set(key, c.Header.Get(key))
	}

	return &Request{c.HttpClient, nil, httpreq.URL.Query(), httpreq}, err
}

func (r *Request) Do(method string) *Response {
	r.URL.RawQuery = r.Query.Encode()
	r.Method = method
	httpres, err := r.Client.Do(r.Request)
	if err != nil {
		return ResponseError(err)
	}

	mtype, err := mediaType(httpres)
	if err != nil {
		httpres.Body.Close()
		return ResponseError(err)
	}

	headerDecoder := mediaheader.Decoder{}
	mheader := headerDecoder.Decode(httpres.Header)

	return &Response{nil, mtype, mheader, UseApiError(httpres.StatusCode), false, httpres}
}

func (r *Request) Head() *Response {
	return r.Do(HeadMethod)
}

func (r *Request) Get() *Response {
	return r.Do(GetMethod)
}

func (r *Request) Post() *Response {
	return r.Do(PostMethod)
}

func (r *Request) Put() *Response {
	return r.Do(PutMethod)
}

func (r *Request) Patch() *Response {
	return r.Do(PatchMethod)
}

func (r *Request) Delete() *Response {
	return r.Do(DeleteMethod)
}

func (r *Request) Options() *Response {
	return r.Do(OptionsMethod)
}

func (r *Request) SetBody(mtype *mediatype.MediaType, input interface{}) error {
	r.MediaType = mtype
	buf, err := mtype.Encode(input)
	if err != nil {
		return err
	}

	r.Header.Set(ctypeHeader, mtype.String())
	r.ContentLength = int64(buf.Len())
	r.Body = ioutil.NopCloser(buf)
	return nil
}

const (
	ctypeHeader   = "Content-Type"
	HeadMethod    = "HEAD"
	GetMethod     = "GET"
	PostMethod    = "POST"
	PutMethod     = "PUT"
	PatchMethod   = "PATCH"
	DeleteMethod  = "DELETE"
	OptionsMethod = "OPTIONS"
)
