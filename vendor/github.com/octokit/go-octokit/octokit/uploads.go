package octokit

import (
	"io"
	"net/url"
)

// Uploads creates an UploadsService with a base url
func (c *Client) Uploads(url *url.URL) (uploads *UploadsService) {
	uploads = &UploadsService{client: c, URL: url}
	return
}

// UploadsService is a service providing access to asset uploads from a particular url
type UploadsService struct {
	client *Client
	URL    *url.URL
}

// UploadAsset uploads a particular asset of some content type and length to the service
func (u *UploadsService) UploadAsset(asset io.ReadCloser, contentType string, contentLength int64) (result *Result) {
	return u.client.upload(u.URL, asset, contentType, contentLength)
}
