package update

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/kr/binarydist"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"testing"
)

var (
	oldFile = []byte{0xDE, 0xAD, 0xBE, 0xEF}
	newFile = []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
)

func cleanup(path string) {
	os.Remove(path)
}

// we write with a separate name for each test so that we can run them in parallel
func writeOldFile(path string, t *testing.T) {
	if err := ioutil.WriteFile(path, oldFile, 0777); err != nil {
		t.Fatalf("Failed to write file for testing preparation: %v", err)
	}
}

func validateUpdate(path string, err error, t *testing.T) {
	if err != nil {
		t.Fatalf("Failed to update: %v", err)
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file post-update: %v", err)
	}

	if !bytes.Equal(buf, newFile) {
		t.Fatalf("File was not updated! Bytes read: %v, Bytes expected: %v", buf, newFile)
	}
}

func TestFromStream(t *testing.T) {
	t.Parallel()

	fName := "TestFromStream"
	defer cleanup(fName)
	writeOldFile(fName, t)

	err, _ := New().Target(fName).FromStream(bytes.NewReader(newFile))
	validateUpdate(fName, err, t)
}

func TestFromFile(t *testing.T) {
	t.Parallel()

	fName := "TestFromFile"
	newFName := "NewTestFromFile"
	defer cleanup(fName)
	defer cleanup(newFName)
	writeOldFile(fName, t)

	if err := ioutil.WriteFile(newFName, newFile, 0777); err != nil {
		t.Fatalf("Failed to write file to update from: %v", err)
	}

	err, _ := New().Target(fName).FromFile(newFName)
	validateUpdate(fName, err, t)
}

func TestFromUrl(t *testing.T) {
	t.Parallel()

	fName := "TestFromUrl"
	defer cleanup(fName)
	writeOldFile(fName, t)

	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Couldn't bind listener: %v", err)
	}
	addr := l.Addr().String()

	go http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(newFile)
	}))

	err, _ = New().Target(fName).FromUrl("http://" + addr)
	validateUpdate(fName, err, t)
}

func TestVerifyChecksum(t *testing.T) {
	t.Parallel()

	fName := "TestVerifyChecksum"
	defer cleanup(fName)
	writeOldFile(fName, t)

	checksum, err := ChecksumForBytes(newFile)
	if err != nil {
		t.Fatalf("Failed to compute checksum: %v", err)
	}

	err, _ = New().Target(fName).VerifyChecksum(checksum).FromStream(bytes.NewReader(newFile))
	validateUpdate(fName, err, t)
}

func TestVerifyChecksumNegative(t *testing.T) {
	t.Parallel()

	fName := "TestVerifyChecksumNegative"
	defer cleanup(fName)
	writeOldFile(fName, t)

	badChecksum := []byte{0x0A, 0x0B, 0x0C, 0xFF}
	err, _ := New().Target(fName).VerifyChecksum(badChecksum).FromStream(bytes.NewReader(newFile))
	if err == nil {
		t.Fatalf("Failed to detect bad checksum!")
	}
}

func TestApplyPatch(t *testing.T) {
	t.Parallel()

	fName := "TestApplyPatch"
	defer cleanup(fName)
	writeOldFile(fName, t)

	patch := new(bytes.Buffer)
	err := binarydist.Diff(bytes.NewReader(oldFile), bytes.NewReader(newFile), patch)
	if err != nil {
		t.Fatalf("Failed to create patch: %v", err)
	}

	up := New().Target(fName).ApplyPatch(PATCHTYPE_BSDIFF)
	err, _ = up.FromStream(bytes.NewReader(patch.Bytes()))
	validateUpdate(fName, err, t)
}

func TestCorruptPatch(t *testing.T) {
	t.Parallel()

	fName := "TestCorruptPatch"
	defer cleanup(fName)
	writeOldFile(fName, t)

	badPatch := []byte{0x44, 0x38, 0x86, 0x3c, 0x4f, 0x8d, 0x26, 0x54, 0xb, 0x11, 0xce, 0xfe, 0xc1, 0xc0, 0xf8, 0x31, 0x38, 0xa0, 0x12, 0x1a, 0xa2, 0x57, 0x2a, 0xe1, 0x3a, 0x48, 0x62, 0x40, 0x2b, 0x81, 0x12, 0xb1, 0x21, 0xa5, 0x16, 0xed, 0x73, 0xd6, 0x54, 0x84, 0x29, 0xa6, 0xd6, 0xb2, 0x1b, 0xfb, 0xe6, 0xbe, 0x7b, 0x70}
	up := New().Target(fName).ApplyPatch(PATCHTYPE_BSDIFF)
	err, _ := up.FromStream(bytes.NewReader(badPatch))
	if err == nil {
		t.Fatalf("Failed to detect corrupt patch!")
	}
}

func TestVerifyChecksumPatchNegative(t *testing.T) {
	t.Parallel()

	fName := "TestVerifyChecksumPatchNegative"
	defer cleanup(fName)
	writeOldFile(fName, t)

	checksum, err := ChecksumForBytes(newFile)
	if err != nil {
		t.Fatalf("Failed to compute checksum: %v", err)
	}

	patch := new(bytes.Buffer)
	anotherFile := []byte{0x77, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	err = binarydist.Diff(bytes.NewReader(oldFile), bytes.NewReader(anotherFile), patch)
	if err != nil {
		t.Fatalf("Failed to create patch: %v", err)
	}

	up := New().Target(fName).ApplyPatch(PATCHTYPE_BSDIFF).VerifyChecksum(checksum)
	err, _ = up.FromStream(bytes.NewReader(patch.Bytes()))
	if err == nil {
		t.Fatalf("Failed to detect patch to wrong file!")
	}
}

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxSWmu7trWKAwDFjiCN2D
Tk2jj2sgcr/CMlI4cSSiIOHrXCFxP1I8i9PvQkd4hasXQrLbT5WXKrRGv1HKUKab
b9ead+kD0kxk7i2bFYvKX43oq66IW0mOLTQBO7I9UyT4L7svcMD+HUQ2BqHoaQe4
y20C59dPr9Dpcz8DZkdLsBV6YKF6Ieb3iGk8oRLMWNaUqPa8f1BGgxAkvPHcqDjT
x4xRnjgTRRRlZvRtALHMUkIChgxDOhoEzKpGiqnX7HtMJfrhV6h0PAXNA4h9Kjv5
5fhJ08Rz7mmZmtH5JxTK5XTquo59sihSajR4bSjZbbkQ1uLkeFlY3eli3xdQ7Nrf
fQIDAQAB
-----END PUBLIC KEY-----`

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAxSWmu7trWKAwDFjiCN2DTk2jj2sgcr/CMlI4cSSiIOHrXCFx
P1I8i9PvQkd4hasXQrLbT5WXKrRGv1HKUKabb9ead+kD0kxk7i2bFYvKX43oq66I
W0mOLTQBO7I9UyT4L7svcMD+HUQ2BqHoaQe4y20C59dPr9Dpcz8DZkdLsBV6YKF6
Ieb3iGk8oRLMWNaUqPa8f1BGgxAkvPHcqDjTx4xRnjgTRRRlZvRtALHMUkIChgxD
OhoEzKpGiqnX7HtMJfrhV6h0PAXNA4h9Kjv55fhJ08Rz7mmZmtH5JxTK5XTquo59
sihSajR4bSjZbbkQ1uLkeFlY3eli3xdQ7NrffQIDAQABAoIBAAkN+6RvrTR61voa
Mvd5RQiZpEN4Bht/Fyo8gH8h0Zh1B9xJZOwlmMZLS5fdtHlfLEhR8qSrGDBL61vq
I8KkhEsUufF78EL+YzxVN+Q7cWYGHIOWFokqza7hzpSxUQO6lPOMQ1eIZaNueJTB
Zu07/47ISPPg/bXzgGVcpYlTCPTjUwKjtfyMqvX9AD7fIyYRm6zfE7EHj1J2sBFt
Yz1OGELg6HfJwXfpnPfBvftD0hWGzJ78Bp71fPJe6n5gnqmSqRvrcXNWFnH/yqkN
d6vPIxD6Z3LjvyZpkA7JillLva2L/zcIFhg4HZvQnWd8/PpDnUDonu36hcj4SC5j
W4aVPLkCgYEA4XzNKWxqYcajzFGZeSxlRHupSAl2MT7Cc5085MmE7dd31wK2T8O4
n7N4bkm/rjTbX85NsfWdKtWb6mpp8W3VlLP0rp4a/12OicVOkg4pv9LZDmY0sRlE
YuDJk1FeCZ50UrwTZI3rZ9IhZHhkgVA6uWAs7tYndONkxNHG0pjqs4sCgYEA39MZ
JwMqo3qsPntpgP940cCLflEsjS9hYNO3+Sv8Dq3P0HLVhBYajJnotf8VuU0fsQZG
grmtVn1yThFbMq7X1oY4F0XBA+paSiU18c4YyUnwax2u4sw9U/Q9tmQUZad5+ueT
qriMBwGv+ewO+nQxqvAsMUmemrVzrfwA5Oct+hcCgYAfiyXoNZJsOy2O15twqBVC
j0oPGcO+/9iT89sg5lACNbI+EdMPNYIOVTzzsL1v0VUfAe08h++Enn1BPcG0VHkc
ZFBGXTfJoXzfKQrkw7ZzbzuOGB4m6DH44xlP0oIlNlVvfX/5ASF9VJf3RiBJNsAA
TsP6ZVr/rw/ZuL7nlxy+IQKBgDhL/HOXlE3yOQiuOec8WsNHTs7C1BXe6PtVxVxi
988pYK/pclL6zEq5G5NLSceF4obAMVQIJ9UtUGbabrncyGUo9UrFPLsjYvprSZo8
YHegpVwL50UcYgCP2kXZ/ldjPIcjYDz8lhvdDMor2cidGTEJn9P11HLNWP9V91Ob
4jCZAoGAPNRSC5cC8iP/9j+s2/kdkfWJiNaolPYAUrmrkL6H39PYYZM5tnhaIYJV
Oh9AgABamU0eb3p3vXTISClVgV7ifq1HyZ7BSUhMfaY2Jk/s3sUHCWFxPZe9sgEG
KinIY/373KIkIV/5g4h2v1w330IWcfptxKcY/Er3DJr38f695GE=
-----END RSA PRIVATE KEY-----`

func sign(privatePEM string, source []byte, t *testing.T) []byte {
	block, _ := pem.Decode([]byte(privatePEM))
	if block == nil {
		t.Fatalf("Failed to parse private key PEM")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		t.Fatalf("Failed to parse private key DER")
	}

	checksum, err := ChecksumForBytes(source)
	if err != nil {
		t.Fatalf("Failed to make checksum")
	}

	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, checksum)
	if err != nil {
		t.Fatalf("Failed to sign: %v", sig)
	}

	return sig
}

func TestVerifySignature(t *testing.T) {
	t.Parallel()

	fName := "TestVerifySignature"
	defer cleanup(fName)
	writeOldFile(fName, t)

	up, err := New().Target(fName).VerifySignatureWithPEM([]byte(publicKey))
	if err != nil {
		t.Fatalf("Could not parse public key: %v", err)
	}

	signature := sign(privateKey, newFile, t)
	err, _ = up.VerifySignature(signature).FromStream(bytes.NewReader(newFile))
	validateUpdate(fName, err, t)
}

func TestVerifyFailBadSignature(t *testing.T) {
	t.Parallel()

	fName := "TestVerifyFailBadSignature"
	defer cleanup(fName)
	writeOldFile(fName, t)

	up, err := New().Target(fName).VerifySignatureWithPEM([]byte(publicKey))
	if err != nil {
		t.Fatalf("Could not parse public key: %v", err)
	}

	badSig := []byte{0xFF, 0xEE, 0xDD, 0xCC, 0xBB, 0xAA}
	err, _ = up.VerifySignature(badSig).FromStream(bytes.NewReader(newFile))
	if err == nil {
		t.Fatalf("Did not fail with bad signature")
	}
}

func TestVerifyFailNoSignature(t *testing.T) {
	t.Parallel()

	fName := "TestVerifySignatureWithPEM"
	defer cleanup(fName)
	writeOldFile(fName, t)

	up, err := New().Target(fName).VerifySignatureWithPEM([]byte(publicKey))
	if err != nil {
		t.Fatalf("Could not parse public key: %v", err)
	}

	err, _ = up.VerifySignature([]byte{}).FromStream(bytes.NewReader(newFile))
	if err == nil {
		t.Fatalf("Did not fail with empty signature")
	}
}

const wrongKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEArKqjT+xOFJILe0CX7lKfQy52YwWLF9devYtLeUHTbPOueGLy
6CjrXJBrWIxNBxRd53y4dtgiMqCX6Gmmvuy8HnfbBuJjR2mcdEYo8UDy+aSVBQ6T
/ND7Fd7KSzOruEFFzl2QFnZ/SrW/nsXdGyuF8l+YIwjZJRyV6StZkZ4ydOzOqUk9
FXTeIkhX/Q7/jTETw7L3wxMyLgJAlV3lxDsPkMjxymngtbAIjwEjLsVeU+Prcz2e
Ww34SZQ8qwzAdXieuDPryMwEsCcgc5NAKJFNL8TppYGDXOHI7CRXqHfNiJq2R+kQ
LdxRvmfx8/iu4xM2hBzk4uSDS6TTn2AnWBm+cQIDAQABAoIBAFp//aUwaCRj/9yU
GI3zhEJEIgz4pNTUL3YNgnuFwvlCJ9o1kreYavRTRdBdiSoCxM1GE7FGy3XZsoVA
iwNbNaaKj6RmGD8f3b8b3u3EaxXp66mA4JQMPO5TnZgY9xJWM+5cH9+GMGXKKStg
7ekFwOkuraD/TEElYHWcIRAv6KZbc/YOIa6YDKi+1Gc7u0MeIvwqN7nwaBAoJKUE
ZrJIfYKIViD/ZrCpgWN47C9x8w3ne7iiDrYoYct+0reC9LFlqwVBtDnyVx/q3upW
zzczbNQagu3w0QgprDGhy0ZhDNxuylV3XBWTB+xBrFQgz6rD3LzUPywlbt0N7ZmD
936MVSECgYEA1IElCahF/+hC/OxFgy98DubAUDGmrvxWeZF3bvTseWZQp/gzxVS+
SYumYyd2Ysx5+UjXQlVgR6BbDG13+DpSpZm6+MeWHBAR+KA2qCg009SDFv7l26/d
xMT7lvIWz7ckQDb/+jvhF9HL2llyTN1Zex+n3XBeAMKNrPaubdEBFsUCgYEA0AIO
tZMtzOpioAR1lGbwIguq04msDdrJNaY2TKrLeviJuQUw94fgL+3ULAPsiyxaU/Gv
vln11R7aIp1SJ09T2UoFRbty+6SGRC56+Wh0pn5VnAi7aT6qdkYWhEjhqRHuXosf
PYboXBuMwA0FBUTxWQL/lux2PZgvBkniYh5jI70CgYEAk9KmhhpFX2gdOT3OeRxO
CzufaemwDqfAK97yGwBLg4OV9dJliQ6TNCvt+amY489jxfJSs3UafZjh3TpFKyq/
FS1kb+y+0hSnu7EPdFhLr1N0QUndcb3b4iY48V7EWYgHspfP5y1CPsSVLvXr2eZc
eZaiuhqReavczAXpfsDWJhUCgYEAwmUp2gfyhc+G3IVOXaLWSPseaxP+9/PAl6L+
nCgCgqpEC+YOHUee/SwHXhtMtcR9pnX5CKyKUuLCehcM8C/y7N+AjerhSsw3rwDB
bNVyLydiWrDOdU1bga1+3aI/QwK/AxyB1b5+6ZXVtKZ2SrZj2Aw1UZcr6eSQDhB+
wbQkcwECgYBF13FMA6OOon992t9H3I+4KDgmz6G6mz3bVXSoFWfO1p/yXP04BzJl
jtLFvFVTZdMs2o/wTd4SL6gYjx9mlOWwM8FblmjfiNSUVIyye33fRntEAr1n+FYI
Xhv6aVnNdaGehGIqQxXFoGyiJxG3RYNkSwaTOamxY1V+ceLuO26n2Q==
-----END RSA PRIVATE KEY-----`

func TestVerifyFailWrongSignature(t *testing.T) {
	t.Parallel()

	fName := "TestVerifyFailWrongSignature"
	defer cleanup(fName)
	writeOldFile(fName, t)

	up, err := New().Target(fName).VerifySignatureWithPEM([]byte(publicKey))
	if err != nil {
		t.Fatalf("Could not parse public key: %v", err)
	}

	signature := sign(wrongKey, newFile, t)
	err, _ = up.VerifySignature(signature).FromStream(bytes.NewReader(newFile))
	if err == nil {
		t.Fatalf("Verified an update that was signed by an untrusted key!")
	}
}

func TestSignatureButNoPublicKey(t *testing.T) {
	t.Parallel()

	fName := "TestSignatureButNoPublicKey"
	defer cleanup(fName)
	writeOldFile(fName, t)

	sig := sign(privateKey, newFile, t)
	err, _ := New().Target(fName).VerifySignature(sig).FromStream(bytes.NewReader(newFile))
	if err == nil {
		t.Fatalf("Allowed an update with a signautre verification when no public key was specified!")
	}
}

func TestPublicKeyButNoSignature(t *testing.T) {
	t.Parallel()

	fName := "TestPublicKeyButNoSignature"
	defer cleanup(fName)
	writeOldFile(fName, t)

	up, err := New().Target(fName).VerifySignatureWithPEM([]byte(publicKey))
	if err != nil {
		t.Fatalf("Could not parse public key: %v", err)
	}

	err, _ = up.FromStream(bytes.NewReader(newFile))
	if err == nil {
		t.Fatalf("Allowed an update with no signautre when a public key was specified!")
	}
}
