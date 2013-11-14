package octokit

import (
	"github.com/lostisland/go-sawyer/mediaheader"
	"github.com/lostisland/go-sawyer/mediatype"
	"net/http"
)

type Response struct {
	MediaType   *mediatype.MediaType
	MediaHeader *mediaheader.MediaHeader
	*http.Response
}
