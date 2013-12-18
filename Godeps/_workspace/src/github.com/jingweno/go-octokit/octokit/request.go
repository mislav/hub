package octokit

import (
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/mediatype"
)

type Request struct {
	sawyerReq *sawyer.Request
}

func (r *Request) Head(output interface{}) (resp *Response, err error) {
	resp, err = r.do(sawyer.HeadMethod, nil, output)
	return
}

func (r *Request) Get(output interface{}) (resp *Response, err error) {
	resp, err = r.do(sawyer.GetMethod, nil, output)
	return
}

func (r *Request) Post(input interface{}, output interface{}) (resp *Response, err error) {
	resp, err = r.do(sawyer.PostMethod, input, output)
	return
}

func (r *Request) Put(input interface{}, output interface{}) (resp *Response, err error) {
	resp, err = r.do(sawyer.PutMethod, input, output)
	return
}

func (r *Request) Delete(output interface{}) (resp *Response, err error) {
	resp, err = r.do(sawyer.DeleteMethod, nil, output)
	return
}

func (r *Request) Patch(input interface{}, output interface{}) (resp *Response, err error) {
	resp, err = r.do(sawyer.PatchMethod, input, output)
	return
}

func (r *Request) do(method string, input interface{}, output interface{}) (resp *Response, err error) {
	var sawyerResp *sawyer.Response
	switch method {
	case sawyer.HeadMethod:
		sawyerResp = r.sawyerReq.Head()
	case sawyer.GetMethod:
		sawyerResp = r.sawyerReq.Get()
	case sawyer.PostMethod:
		mtype, _ := mediatype.Parse(defaultMediaType)
		r.sawyerReq.SetBody(mtype, input)
		sawyerResp = r.sawyerReq.Post()
	case sawyer.PutMethod:
		mtype, _ := mediatype.Parse(defaultMediaType)
		r.sawyerReq.SetBody(mtype, input)
		sawyerResp = r.sawyerReq.Put()
	case sawyer.PatchMethod:
		mtype, _ := mediatype.Parse(defaultMediaType)
		r.sawyerReq.SetBody(mtype, input)
		sawyerResp = r.sawyerReq.Patch()
	case sawyer.DeleteMethod:
		sawyerResp = r.sawyerReq.Delete()
	case sawyer.OptionsMethod:
		sawyerResp = r.sawyerReq.Options()
	}

	if sawyerResp.IsError() {
		err = sawyerResp.ResponseError
		return
	}

	if sawyerResp.IsApiError() {
		err = NewResponseError(sawyerResp)
		return
	}

	resp = &Response{Response: sawyerResp.Response, MediaType: sawyerResp.MediaType, MediaHeader: sawyerResp.MediaHeader}
	err = sawyerResp.Decode(output)

	return
}
