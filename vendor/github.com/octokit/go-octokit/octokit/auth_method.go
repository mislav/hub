package octokit

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/fhs/go-netrc/netrc"
)

// AuthMethod is a general interface for possible forms of authentication.
// In order to act as a form of authentication, a struct must at minimum
// implement the String() method which produces the authentication string.
// See http://developer.github.com/v3/auth/ for more information
type AuthMethod interface {
	fmt.Stringer
}

// BasicAuth is a form of authentication involving a simple login and password.
// A OneTimePassword field may be set for two-factor authentication.
type BasicAuth struct {
	Login           string
	Password        string
	OneTimePassword string
}

// String hashes the login and password to produce the string to be passed for
// authentication purposes.
func (b BasicAuth) String() string {
	return fmt.Sprintf("Basic %s", hashAuth(b.Login, b.Password))
}

// NetrcAuth is a form of authentication using a .netrc file for permanent
// authentication by storing credentials.
type NetrcAuth struct {
	NetrcPath string
}

// String accesses the credentials from the .netrc file and hashes the associated
// login and password to submit as a form of basic authentication.
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

// hashAuth is a helper function for producing a base64 encoding of a username and
// password pair
func hashAuth(u, p string) string {
	var a = fmt.Sprintf("%s:%s", u, p)
	return base64.StdEncoding.EncodeToString([]byte(a))
}

// TokenAuth is a form of authentication using an access token
type TokenAuth struct {
	AccessToken string
}

// String produces the authentication string using the access token
func (t TokenAuth) String() string {
	return fmt.Sprintf("token %s", t.AccessToken)
}
