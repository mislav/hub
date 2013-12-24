package octokit

import (
	"net/url"
	"os"
)

// Create an UploadsService with the base url.URL
func (c *Client) Uploads(url *url.URL) *UploadsService {
	return &UploadsService{client: c, URL: url}
}

type UploadsService struct {
	client *Client
	URL    *url.URL
}

func (u *UploadsService) UploadAsset(asset *os.File, contentType string) (result *Result) {
	return u.client.upload(u.URL, asset, contentType)
}
