package octokit

import (
	"testing"

	"github.com/jingweno/go-sawyer/hypermedia"
	"github.com/jingweno/go-sawyer/mediaheader"
	"github.com/stretchr/testify/assert"
)

func TestNewResult_Pageable(t *testing.T) {
	resp := &Response{MediaHeader: &mediaheader.MediaHeader{Relations: hypermedia.Relations{"next": hypermedia.Hyperlink("/path")}}}
	result := newResult(resp, nil)

	assert.Equal(t, "/path", string(*result.NextPage))
	assert.Nil(t, result.PrevPage)
	assert.Nil(t, result.LastPage)
	assert.Nil(t, result.FirstPage)
}
