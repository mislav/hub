package check

import (
	"bytes"
	_ "crypto/sha512" // for tls cipher support
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"

	"bitbucket.org/kardianos/osext"
	"github.com/inconshreveable/go-update"
)

type Initiative string

const (
	INITIATIVE_NEVER  Initiative = "never"
	INITIATIVE_AUTO              = "auto"
	INITIATIVE_MANUAL            = "manual"
)

var NoUpdateAvailable error = fmt.Errorf("No update available")

type Params struct {
	// protocol version
	Version int `json:"version"`
	// identifier of the application to update
	AppId string `json:"app_id"`
	// version of the application updating itself
	AppVersion string `json:"app_version"`
	// operating system of target platform
	OS string `json:"-"`
	// hardware architecture of target platform
	Arch string `json:"-"`
	// application-level user identifier
	UserId string `json:"user_id"`
	// checksum of the binary to replace (used for returning diff patches)
	Checksum string `json:"checksum"`
	// release channel (empty string means 'stable')
	Channel string `json:"-"`
	// tags for custom update channels
	Tags map[string]string `json:"tags"`
}

type Result struct {
	up *update.Update

	// should the update be applied automatically/manually
	Initiative Initiative `json:"initiative"`
	// url where to download the updated application
	Url string `json:"url"`
	// a URL to a patch to apply
	PatchUrl string `json:"patch_url"`
	// the patch format (only bsdiff supported at the moment)
	PatchType update.PatchType `json:"patch_type"`
	// version of the new application
	Version string `json:"version"`
	// expected checksum of the new application
	Checksum string `json:"checksum"`
	// signature for verifying update authenticity
	Signature string `json:"signature"`
}

// CheckForUpdate makes an HTTP post to a URL with the JSON serialized
// representation of Params. It returns the deserialized result object
// returned by the remote endpoint or an error. If you do not set
// OS/Arch, CheckForUpdate will populate them for you. Similarly, if
// Version is 0, it will be set to 1. Lastly, if Checksum is the empty
// string, it will be automatically be computed for the running program's
// executable file.
func (p *Params) CheckForUpdate(url string, up *update.Update) (*Result, error) {
	if p.Tags == nil {
		p.Tags = make(map[string]string)
	}

	if p.Channel == "" {
		p.Channel = "stable"
	}

	if p.OS == "" {
		p.OS = runtime.GOOS
	}

	if p.Arch == "" {
		p.Arch = runtime.GOARCH
	}

	if p.Version == 0 {
		p.Version = 1
	}

	// ignore errors auto-populating the checksum
	// if it fails, you just won't be able to patch
	if up.TargetPath == "" {
		p.Checksum = defaultChecksum()
	} else {
		checksum, err := update.ChecksumForFile(up.TargetPath)
		if err != nil {
			return nil, err
		}
		p.Checksum = hex.EncodeToString(checksum)
	}

	p.Tags["os"] = p.OS
	p.Tags["arch"] = p.Arch
	p.Tags["channel"] = p.Channel

	body, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// no content means no available update
	if resp.StatusCode == 204 {
		return nil, NoUpdateAvailable
	}

	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &Result{up: up}
	if err := json.Unmarshal(respBytes, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *Params) CheckAndApplyUpdate(url string, up *update.Update) (result *Result, err error, errRecover error) {
	// check for an update
	result, err = p.CheckForUpdate(url, up)
	if err != nil {
		return
	}

	// run the available update
	err, errRecover = result.Update()
	return
}

func (r *Result) Update() (err error, errRecover error) {
	if r.Checksum != "" {
		r.up.Checksum, err = hex.DecodeString(r.Checksum)
		if err != nil {
			return
		}
	}

	if r.Signature != "" {
		r.up.Signature, err = hex.DecodeString(r.Signature)
		if err != nil {
			return
		}
	}

	if r.PatchType != "" {
		r.up.PatchType = r.PatchType
	}

	if r.Url == "" && r.PatchUrl == "" {
		err = fmt.Errorf("Result does not contain an update url or patch update url")
		return
	}

	if r.PatchUrl != "" {
		err, errRecover = r.up.FromUrl(r.PatchUrl)
		if err == nil {
			// success!
			return
		} else {
			// failed to update from patch URL, try with the whole thing
			if r.Url == "" || errRecover != nil {
				// we can't try updating from a URL with the full contents
				// in these cases, so fail
				return
			} else {
				r.up.PatchType = update.PATCHTYPE_NONE
			}
		}
	}

	// try updating from a URL with the full contents
	return r.up.FromUrl(r.Url)
}

func defaultChecksum() string {
	path, err := osext.Executable()
	if err != nil {
		return ""
	}

	checksum, err := update.ChecksumForFile(path)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(checksum)
}
