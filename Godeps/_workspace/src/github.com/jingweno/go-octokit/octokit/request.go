package octokit

import (
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/mediatype"
)

type Request struct {
	sawyerReq *sawyer.Request
}

func (r *Request) Head(output interface{}) (*Response, error) {
	return r.createResponse(r.sawyerReq.Head(), output)
}

func (r *Request) Get(output interface{}) (*Response, error) {
	return r.createResponse(r.sawyerReq.Get(), output)
}

func (r *Request) Post(input interface{}, output interface{}) (*Response, error) {
	r.setBody(input)
	return r.createResponse(r.sawyerReq.Post(), output)
}

func (r *Request) Put(input interface{}, output interface{}) (*Response, error) {
	r.setBody(input)
	return r.createResponse(r.sawyerReq.Put(), output)
}

func (r *Request) Delete(output interface{}) (*Response, error) {
	return r.createResponse(r.sawyerReq.Delete(), output)
}

func (r *Request) Patch(input interface{}, output interface{}) (*Response, error) {
	r.setBody(input)
	return r.createResponse(r.sawyerReq.Patch(), output)
}

func (r *Request) Options(output interface{}) (*Response, error) {
	return r.createResponse(r.sawyerReq.Options(), output)
}

func (r *Request) setBody(input interface{}) {
	mtype, _ := mediatype.Parse(defaultMediaType)
	r.sawyerReq.SetBody(mtype, input)
}

func (r *Request) createResponse(sawyerResp *sawyer.Response, output interface{}) (resp *Response, err error) {
	resp, err = NewResponse(sawyerResp)
	if err == nil {
		err = sawyerResp.Decode(output)
	}

	return
}
