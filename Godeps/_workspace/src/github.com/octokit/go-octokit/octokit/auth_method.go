package octokit

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/fhs/go-netrc/netrc"
)

// See http://developer.github.com/v3/auth/
type AuthMethod interface {
	fmt.Stringer
}

type BasicAuth struct {
	Login           string
	Password        string
	OneTimePassword string // for two-factor authentication
}

func (b BasicAuth) String() string {
	return fmt.Sprintf("Basic %s", hashAuth(b.Login, b.Password))
}

type NetrcAuth struct {
	NetrcPath string
}

func (n NetrcAuth) String() string {
	netrcPath := n.NetrcPath
	if netrcPath == "" {
		netrcPath = filepath.Join(os.Getenv("HOME"), ".netrc")
	}
	apiURL, _ := url.Parse(gitHubAPIURL)
	credentials, err := netrc.FindMachine(netrcPath, apiURL.Host)
	if err != nil {
		panic(fmt.Errorf("netrc error (%s): %v", apiURL.Host, err))
	}
	return fmt.Sprintf("Basic %s", hashAuth(credentials.Login, credentials.Password))
}

func hashAuth(u, p string) string {
	var a = fmt.Sprintf("%s:%s", u, p)
	return base64.StdEncoding.EncodeToString([]byte(a))
}

type TokenAuth struct {
	AccessToken string
}

func (t TokenAuth) String() string {
	return fmt.Sprintf("token %s", t.AccessToken)
}
