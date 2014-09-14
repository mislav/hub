/*
go-update allows a program to update itself by replacing its executable file
with a new version. It provides the flexibility to implement different updating user experiences
like auto-updating, or manual user-initiated updates. It also boasts
advanced features like binary patching and code signing verification.

Updating your program to a new version is as easy as:

	err, errRecover := update.New().FromUrl("http://release.example.com/2.0/myprogram")
	if err != nil {
		fmt.Printf("Update failed: %v\n", err)
	}

You may also choose to update from other data sources such as a file or an io.Reader:

    err, errRecover := update.New().FromFile("/path/to/update")

Binary Diff Patching

Binary diff updates are supported and easy to use:

    up := update.New().ApplyPatch(update.PATCHTYPE_BSDIFF)
    err, errRecover := up.FromUrl("http://release.example.com/2.0/mypatch")

Checksum Verification

You should also verify the checksum of new updates as well as verify
the digital signature of an update. Note that even when you choose to apply
a patch, the checksum is verified against the complete update after that patch
has been applied.

    up := update.New().ApplyPatch(update.PATCHTYPE_BSDIFF).VerifyChecksum(checksum)
    err, errRecover := up.FromUrl("http://release.example.com/2.0/mypatch")

Updating other files

Updating arbitrary files is also supported. You may update files which are
not the currently running program:

    up := update.New().Target("/usr/local/bin/some-program")
    err, errRecover := up.FromUrl("http://release.example.com/2.0/some-program")

Code Signing

Truly secure updates use code signing to verify that the update was issued by a trusted party.
To do this, you'll need to generate a public/private key pair. You can do this with openssl,
or the equinox.io client (https://equinox.io/client) can easily generate one for you:

    # with equinox client
    equinox genkey --private-key=private.pem --public-key=public.pem

    # with openssl
    openssl genrsa -out private.pem 2048
    openssl rsa -in private.pem -out public.pem -pubout

Once you have your key pair, you can instruct your program to validate its updates
with the public key:

    const publicKey = `-----BEGIN PUBLIC KEY-----
    ...
    -----END PUBLIC KEY-----`

    up, err := update.New().VerifySignatureWithPEM(publicKey)
    if err != nil {
        return fmt.Errorf("Bad public key: '%v': %v", publicKey, err)
    }

Once you've configured your program this way, it will disallow all updates unless they
are properly signed. You must now pass in the signature to verify with:

    up.VerifySignature(signature).FromUrl("http://dl.example.com/update")

Error Handling and Recovery

To perform an update, the process must be able to read its executable file and to write
to the directory that contains its executable file. It can be useful to check whether the process
has the necessary permissions to perform an update before trying to apply one. Use the
CanUpdate call to provide a useful message to the user if the update can't proceed without
elevated permissions:

    up := update.New().Target("/etc/hosts")
    err := up.CanUpdate()
    if err != nil {
        fmt.Printf("Can't update because: '%v'. Try as root or Administrator\n", err)
        return
    }
    err, errRecover := up.FromUrl("https://example.com/new/hosts")

Although exceedingly unlikely, the update operation itself is not atomic and can fail
in such a way that a user's computer is left in an inconsistent state. If that happens,
go-update attempts to recover to leave the system in a good state. If the recovery step
fails (even more unlikely), a second error, referred to as "errRecover" will be non-nil
so that you may inform your users of the bad news. You should handle this case as shown
here:

    err, errRecover := up.FromUrl("https://example.com/update")
    if err != nil {
        fmt.Printf("Update failed: %v\n", err)
        if errRecover != nil {
            fmt.Printf("Failed to recover bad update: %v!\n", errRecover)
            fmt.Printf("Program exectuable may be missing!\n")
        }
    }

Subpackages

Sub-package check contains the client functionality for a simple protocol for negotiating
whether a new update is available, where it is, and the metadata needed for verifying it.

Sub-package download contains functionality for downloading from an HTTP endpoint
while outputting a progress meter and supports resuming partial downloads.
*/
package update

import (
	"bitbucket.org/kardianos/osext"
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	_ "crypto/sha512" // for tls cipher support
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/inconshreveable/go-update/download"
	"github.com/kr/binarydist"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// The type of a binary patch, if any. Only bsdiff is supported
type PatchType string

const (
	PATCHTYPE_BSDIFF PatchType = "bsdiff"
	PATCHTYPE_NONE             = ""
)

type Update struct {
	// empty string means "path of the current executable"
	TargetPath string

	// type of patch to apply. PATCHTYPE_NONE means "not a patch"
	PatchType

	// sha256 checksum of the new binary to verify against
	Checksum []byte

	// public key to use for signature verification
	PublicKey *rsa.PublicKey

	// signature to use for signature verification
	Signature []byte
}

func (u *Update) getPath() (string, error) {
	if u.TargetPath == "" {
		return osext.Executable()
	} else {
		return u.TargetPath, nil
	}
}

// New creates a new Update object.
// A default update object assumes the complete binary
// content will be used for update (not a patch) and that
// the intended target is the running executable.
//
// Use this as the start of a chain of calls on the Update
// object to build up your configuration. Example:
//
//     up := update.New().ApplyPatch(update.PATCHTYPE_BSDIFF).VerifyChecksum(checksum)
//
func New() *Update {
	return &Update{
		TargetPath: "",
		PatchType:  PATCHTYPE_NONE,
	}
}

// Target configures the update to update the file at the given path.
// The emptry string means 'the executable file of the running program'.
func (u *Update) Target(path string) *Update {
	u.TargetPath = path
	return u
}

// ApplyPatch configures the update to treat the contents of the update
// as a patch to apply to the existing to target. You must specify the
// format of the patch. Only PATCHTYPE_BSDIFF is supported at the moment.
func (u *Update) ApplyPatch(patchType PatchType) *Update {
	u.PatchType = patchType
	return u
}

// VerifyChecksum configures the update to verify that the
// the update has the given sha256 checksum.
func (u *Update) VerifyChecksum(checksum []byte) *Update {
	u.Checksum = checksum
	return u
}

// VerifySignature configures the update to verify the given
// signature of the update. You must also call one of the
// VerifySignatureWith* functions to specify a public key
// to use for verification.
func (u *Update) VerifySignature(signature []byte) *Update {
	u.Signature = signature
	return u
}

// VerifySignatureWith configures the update to use the given RSA
// public key to verify the update's signature. You must also call
// VerifySignature() with a signature to check.
//
// You'll probably want to use VerifySignatureWithPEM instead of
// parsing the public key yourself.
func (u *Update) VerifySignatureWith(publicKey *rsa.PublicKey) *Update {
	u.PublicKey = publicKey
	return u
}

// VerifySignatureWithPEM configures the update to use the given PEM-formatted
// RSA public key to verify the update's signature. You must also call
// VerifySignature() with a signature to check.
//
// A PEM formatted public key typically begins with
//     -----BEGIN PUBLIC KEY-----
func (u *Update) VerifySignatureWithPEM(publicKeyPEM []byte) (*Update, error) {
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return u, fmt.Errorf("Couldn't parse PEM data")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return u, err
	}

	var ok bool
	u.PublicKey, ok = pub.(*rsa.PublicKey)
	if !ok {
		return u, fmt.Errorf("Public key isn't an RSA public key")
	}

	return u, nil
}

// FromUrl updates the target with the contents of the given URL.
func (u *Update) FromUrl(url string) (err error, errRecover error) {
	target := new(download.MemoryTarget)
	err = download.New(url, target).Get()
	if err != nil {
		return
	}

	return u.FromStream(target)
}

// FromFile updates the target the contents of the given file.
func (u *Update) FromFile(path string) (err error, errRecover error) {
	// open the new updated contents
	fp, err := os.Open(path)
	if err != nil {
		return
	}
	defer fp.Close()

	// do the update
	return u.FromStream(fp)
}

// FromStream updates the target file with the contents of the supplied io.Reader.
//
// FromStream performs the following actions to ensure a safe cross-platform update:
//
// 1. If configured, applies the contents of the io.Reader as a binary patch.
//
// 2. If configured, computes the sha256 checksum and verifies it matches.
//
// 3. If configured, verifies the RSA signature with a public key.
//
// 4. Creates a new file, /path/to/.target.new with mode 0755 with the contents of the updated file
//
// 5. Renames /path/to/target to /path/to/.target.old
//
// 6. Renames /path/to/.target.new to /path/to/target
//
// 7. If the rename is successful, deletes /path/to/.target.old, returns no error
//
// 8. If the rename fails, attempts to rename /path/to/.target.old back to /path/to/target
// If this operation fails, it is reported in the errRecover return value so as not to
// mask the original error that caused the recovery attempt.
//
// On Windows, the removal of /path/to/.target.old always fails, so instead,
// we just make the old file hidden instead.
func (u *Update) FromStream(updateWith io.Reader) (err error, errRecover error) {
	updatePath, err := u.getPath()
	if err != nil {
		return
	}

	var newBytes []byte
	// apply a patch if requested
	switch u.PatchType {
	case PATCHTYPE_BSDIFF:
		newBytes, err = applyPatch(updateWith, updatePath)
		if err != nil {
			return
		}
	case PATCHTYPE_NONE:
		// no patch to apply, go on through
		newBytes, err = ioutil.ReadAll(updateWith)
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("Unrecognized patch type: %s", u.PatchType)
		return
	}

	// verify checksum if requested
	if u.Checksum != nil {
		if err = verifyChecksum(newBytes, u.Checksum); err != nil {
			return
		}
	}

	// verify signature if requested
	if u.Signature != nil || u.PublicKey != nil {
		if u.Signature == nil {
			err = fmt.Errorf("No public key specified to verify signature")
			return
		}

		if u.PublicKey == nil {
			err = fmt.Errorf("No signature to verify!")
			return
		}

		if err = verifySignature(newBytes, u.Signature, u.PublicKey); err != nil {
			return
		}
	}

	// get the directory the executable exists in
	updateDir := filepath.Dir(updatePath)
	filename := filepath.Base(updatePath)

	// Copy the contents of of newbinary to a the new executable file
	newPath := filepath.Join(updateDir, fmt.Sprintf(".%s.new", filename))
	fp, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return
	}
	defer fp.Close()
	_, err = io.Copy(fp, bytes.NewReader(newBytes))

	// if we don't call fp.Close(), windows won't let us move the new executable
	// because the file will still be "in use"
	fp.Close()

	// this is where we'll move the executable to so that we can swap in the updated replacement
	oldPath := filepath.Join(updateDir, fmt.Sprintf(".%s.old", filename))

	// delete any existing old exec file - this is necessary on Windows for two reasons:
	// 1. after a successful update, Windows can't remove the .old file because the process is still running
	// 2. windows rename operations fail if the destination file already exists
	_ = os.Remove(oldPath)

	// move the existing executable to a new file in the same directory
	err = os.Rename(updatePath, oldPath)
	if err != nil {
		return
	}

	// move the new exectuable in to become the new program
	err = os.Rename(newPath, updatePath)

	if err != nil {
		// copy unsuccessful
		errRecover = os.Rename(oldPath, updatePath)
	} else {
		// copy successful, remove the old binary
		errRemove := os.Remove(oldPath)

		// windows has trouble with removing old binaries, so hide it instead
		if errRemove != nil {
			_ = hideFile(oldPath)
		}
	}

	return
}

// CanUpdate() determines whether the process has the correct permissions to
// perform the requested update. If the update can proceed, it returns nil, otherwise
// it returns the error that would occur if an update were attempted.
func (u *Update) CanUpdate() (err error) {
	// get the directory the file exists in
	path, err := u.getPath()
	if err != nil {
		return
	}

	fileDir := filepath.Dir(path)
	fileName := filepath.Base(path)

	// attempt to open a file in the file's directory
	newPath := filepath.Join(fileDir, fmt.Sprintf(".%s.new", fileName))
	fp, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return
	}
	fp.Close()

	_ = os.Remove(newPath)
	return
}

func applyPatch(patch io.Reader, updatePath string) ([]byte, error) {
	// open the file to update
	old, err := os.Open(updatePath)
	if err != nil {
		return nil, err
	}
	defer old.Close()

	// apply the patch
	applied := new(bytes.Buffer)
	if err = binarydist.Patch(old, applied, patch); err != nil {
		return nil, err
	}

	return applied.Bytes(), nil
}

func verifyChecksum(updated []byte, expectedChecksum []byte) error {
	checksum, err := ChecksumForBytes(updated)
	if err != nil {
		return err
	}

	if !bytes.Equal(expectedChecksum, checksum) {
		return fmt.Errorf("Updated file has wrong checksum. Expected: %x, got: %x", expectedChecksum, checksum)
	}

	return nil
}

// ChecksumForFile returns the sha256 checksum for the given file
func ChecksumForFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ChecksumForReader(f)
}

// ChecksumForReader returns the sha256 checksum for the entire
// contents of the given reader.
func ChecksumForReader(rd io.Reader) ([]byte, error) {
	h := sha256.New()
	if _, err := io.Copy(h, rd); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// ChecksumForBytes returns the sha256 checksum for the given bytes
func ChecksumForBytes(source []byte) ([]byte, error) {
	return ChecksumForReader(bytes.NewReader(source))
}

func verifySignature(source, signature []byte, publicKey *rsa.PublicKey) error {
	checksum, err := ChecksumForBytes(source)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, checksum, signature)
}
