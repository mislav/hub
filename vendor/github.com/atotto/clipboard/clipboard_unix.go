// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build freebsd linux netbsd openbsd solaris dragonfly

package clipboard

import (
	"errors"
	"os/exec"
)

const (
	xsel  = "xsel"
	xclip = "xclip"
)

var (
	Primary bool

	pasteCmdArgs []string
	copyCmdArgs  []string

	xselPasteArgs = []string{xsel, "--output", "--clipboard"}
	xselCopyArgs  = []string{xsel, "--input", "--clipboard"}

	xclipPasteArgs = []string{xclip, "-out", "-selection", "clipboard"}
	xclipCopyArgs  = []string{xclip, "-in", "-selection", "clipboard"}

	missingCommands = errors.New("No clipboard utilities available. Please install xsel or xclip.")
)

func init() {
	pasteCmdArgs = xclipPasteArgs
	copyCmdArgs = xclipCopyArgs

	if _, err := exec.LookPath(xclip); err == nil {
		return
	}

	pasteCmdArgs = xselPasteArgs
	copyCmdArgs = xselCopyArgs

	if _, err := exec.LookPath(xsel); err == nil {
		return
	}

	Unsupported = true
}

func getPasteCommand() *exec.Cmd {
	if Primary {
		pasteCmdArgs = pasteCmdArgs[:1]
	}
	return exec.Command(pasteCmdArgs[0], pasteCmdArgs[1:]...)
}

func getCopyCommand() *exec.Cmd {
	if Primary {
		copyCmdArgs = copyCmdArgs[:1]
	}
	return exec.Command(copyCmdArgs[0], copyCmdArgs[1:]...)
}

func readAll() (string, error) {
	if Unsupported {
		return "", missingCommands
	}
	pasteCmd := getPasteCommand()
	out, err := pasteCmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func writeAll(text string) error {
	if Unsupported {
		return missingCommands
	}
	copyCmd := getCopyCommand()
	in, err := copyCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := copyCmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write([]byte(text)); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}
	return copyCmd.Wait()
}
