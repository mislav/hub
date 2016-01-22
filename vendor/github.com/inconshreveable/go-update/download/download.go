package download

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
)

type roundTripper struct {
	RoundTripFn func(*http.Request) (*http.Response, error)
}

func (rt *roundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt.RoundTripFn(r)
}

// Download encapsulates the state and parameters to download content
// from a URL which:
//
// - Publishes the percentage of the download completed to a channel.
// - May resume a previous download that was partially completed.
//
// Create an instance with the New() factory function.
type Download struct {
	// net/http.Client to use when downloading the update.
	// If nil, a default http.Client is used
	HttpClient *http.Client

	// As bytes are downloaded, they are written to Target.
	// Download also uses the Target's Seek method to determine
	// the size of partial-downloads so that it may properly
	// request the remaining bytes to resume the download.
	Target Target

	// Progress returns the percentage of the download
	// completed as an integer between 0 and 100
	Progress chan (int)

	// HTTP Method to use in the download request. Default is "GET"
	Method string

	// HTTP URL to issue the download request to
	Url string
}

// New initializes a new Download object which will download
// the content from url into target.
func New(url string, target Target) *Download {
	return &Download{
		HttpClient: new(http.Client),
		Progress:   make(chan int),
		Method:     "GET",
		Url:        url,
		Target:     target,
	}
}

// Get() downloads the content of a url to a target destination.
//
// Only HTTP/1.1 servers that implement the Range header support resuming a
// partially completed download.
//
// On success, the server must return 200 and the content, or 206 when resuming a partial download.
// If the HTTP server returns a 3XX redirect, it will be followed according to d.HttpClient's redirect policy.
//
func (d *Download) Get() (err error) {
	// Close the progress channel whenever this function completes
	defer close(d.Progress)

	// determine the size of the download target to determine if we're resuming a partial download
	offset, err := d.Target.Size()
	if err != nil {
		return
	}

	// create the download request
	req, err := http.NewRequest(d.Method, d.Url, nil)
	if err != nil {
		return
	}

	// we have to add headers like this so they get used across redirects
	trans := d.HttpClient.Transport
	if trans == nil {
		trans = http.DefaultTransport
	}

	d.HttpClient.Transport = &roundTripper{
		RoundTripFn: func(r *http.Request) (*http.Response, error) {
			// add header for download continuation
			if offset > 0 {
				r.Header.Add("Range", fmt.Sprintf("%d-", offset))
			}

			// ask for gzipped content so that net/http won't unzip it for us
			// and destroy the content length header we need for progress calculations
			r.Header.Add("Accept-Encoding", "gzip")

			return trans.RoundTrip(r)
		},
	}

	// issue the download request
	resp, err := d.HttpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	// ok
	case 200, 206:

	// server error
	default:
		err = fmt.Errorf("Non 2XX response when downloading update: %s", resp.Status)
		return
	}

	// Determine how much we have to download
	// net/http sets this to -1 when it is unknown
	clength := resp.ContentLength

	// Read the content from the response body
	rd := resp.Body

	// meter the rate at which we download content for
	// progress reporting if we know how much to expect
	if clength > 0 {
		rd = &meteredReader{rd: rd, totalSize: clength, progress: d.Progress}
	}

	// Decompress the content if necessary
	if resp.Header.Get("Content-Encoding") == "gzip" {
		rd, err = gzip.NewReader(rd)
		if err != nil {
			return
		}
	}

	// Download the update
	_, err = io.Copy(d.Target, rd)
	if err != nil {
		return
	}

	return
}

// meteredReader wraps a ReadCloser. Calls to a meteredReader's Read() method
// publish updates to a progress channel with the percentage read so far.
type meteredReader struct {
	rd        io.ReadCloser
	totalSize int64
	progress  chan int
	totalRead int64
	ticks     int64
}

func (m *meteredReader) Close() error {
	return m.rd.Close()
}

func (m *meteredReader) Read(b []byte) (n int, err error) {
	chunkSize := (m.totalSize / 100) + 1
	lenB := int64(len(b))

	var nChunk int
	for start := int64(0); start < lenB; start += int64(nChunk) {
		end := start + chunkSize
		if end > lenB {
			end = lenB
		}

		nChunk, err = m.rd.Read(b[start:end])

		n += nChunk
		m.totalRead += int64(nChunk)

		if m.totalRead > (m.ticks * chunkSize) {
			m.ticks += 1
			// try to send on channel, but don't block if it's full
			select {
			case m.progress <- int(m.ticks + 1):
			default:
			}

			// give the progress channel consumer a chance to run
			runtime.Gosched()
		}

		if err != nil {
			return
		}
	}

	return
}

// A Target is what you can supply to Download,
// it's just an io.Writer with a Size() method so that
// the a Download can "resume" an interrupted download
type Target interface {
	io.Writer
	Size() (int, error)
}

type FileTarget struct {
	*os.File
}

func (t *FileTarget) Size() (int, error) {
	if fi, err := t.File.Stat(); err != nil {
		return 0, err
	} else {
		return int(fi.Size()), nil
	}
}

type MemoryTarget struct {
	bytes.Buffer
}

func (t *MemoryTarget) Size() (int, error) {
	return t.Buffer.Len(), nil
}
