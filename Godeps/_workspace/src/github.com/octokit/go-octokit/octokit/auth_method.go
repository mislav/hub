package octokit

import (
	"encoding/base64"
	"fmt"
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
