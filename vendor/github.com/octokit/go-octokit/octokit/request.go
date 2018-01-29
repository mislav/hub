package octokit

import (
	"github.com/jingweno/go-sawyer"
	"github.com/jingweno/go-sawyer/mediatype"
)

func newRequest(client *Client, urlStr string) (req *Request, err error) {
	sawyerReq, err := client.Client.NewRequest(urlStr)
	if err != nil {
		return
	}

	req = &Request{client: client, Request: sawyerReq}

	return
}

// Request wraps a sawyer Request which is a wrapper for an HttpRequest with
// a particular octokit Client
type Request struct {
	*sawyer.Request
	client *Client
}

// Head sends a HEAD request through the given client and returns the response
// and any associated errors
func (r *Request) Head(output interface{}) (*Response, error) {
	return r.createResponse(r.Request.Head(), output)
}

// Get sends a GET request through the given client and returns the response
// and any associated errors
func (r *Request) Get(output interface{}) (*Response, error) {
	if output == nil {
		return NewResponse(r.Request.Get())
	}
	return r.createResponse(r.Request.Get(), output)
}

// Post sends a POST request through the given client and returns the response
// and any associated errors
func (r *Request) Post(input interface{}, output interface{}) (*Response, error) {
	r.setBody(input)
	return r.createResponse(r.Request.Post(), output)
}

// Put sends a PUT request through the given client and returns the response
// and any associated errors
func (r *Request) Put(input interface{}, output interface{}) (*Response, error) {
	r.setBody(input)
	return r.createResponse(r.Request.Put(), output)
}

// Delete sends a DELETE request through the given client and returns the response
// and any associated errors
func (r *Request) Delete(input interface{}, output interface{}) (*Response, error) {
	r.setBody(input)
	return r.createResponse(r.Request.Delete(), output)
}

// Patch sends a PATCH request through the given client and returns the response
// and any associated errors
func (r *Request) Patch(input interface{}, output interface{}) (*Response, error) {
	r.setBody(input)
	return r.createResponse(r.Request.Patch(), output)
}

// Options sends an OPTIONS request through the given client and returns the response
// and any associated errors
func (r *Request) Options(output interface{}) (*Response, error) {
	return r.createResponse(r.Request.Options(), output)
}

func (r *Request) setBody(input interface{}) {
	mtype, _ := mediatype.Parse(defaultMediaType)
	r.Request.SetBody(mtype, input)
}

func (r *Request) createResponse(sawyerResp *sawyer.Response, output interface{}) (resp *Response, err error) {
	resp, err = NewResponse(sawyerResp)
	if err == nil {
		err = sawyerResp.Decode(output)
	}

	return
}
