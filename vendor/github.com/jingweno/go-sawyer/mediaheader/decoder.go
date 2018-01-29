package mediaheader

import (
	"github.com/jingweno/go-sawyer/hypermedia"
	"net/http"
	"net/url"
	"strings"
)

// TODO: need a full link header parser for http://tools.ietf.org/html/rfc5988
type Decoder struct {
}

func (d *Decoder) Decode(header http.Header) (mediaHeader *MediaHeader) {
	mediaHeader = &MediaHeader{Relations: hypermedia.Relations{}}

	link := header.Get("Link")
	if len(link) == 0 {
		return
	}

	for _, l := range strings.Split(link, ",") {
		l = strings.TrimSpace(l)
		segments := strings.Split(l, ";")

		if len(segments) < 2 {
			continue
		}

		if !strings.HasPrefix(segments[0], "<") || !strings.HasSuffix(segments[0], ">") {
			continue
		}

		url, err := url.Parse(segments[0][1 : len(segments[0])-1])
		if err != nil {
			continue
		}

		link := hypermedia.Hyperlink(url.String())

		for _, segment := range segments[1:] {
			switch strings.TrimSpace(segment) {
			case `rel="next"`:
				mediaHeader.Relations["next"] = link
			case `rel="prev"`:
				mediaHeader.Relations["prev"] = link
			case `rel="first"`:
				mediaHeader.Relations["first"] = link
			case `rel="last"`:
				mediaHeader.Relations["last"] = link
			}
		}
	}

	return
}
