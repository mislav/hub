// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package clipboard

import (
	"syscall"
	"unsafe"
)

const (
	cfUnicodetext = 13
	gmemFixed     = 0x0000
)

var (
	user32           = syscall.MustLoadDLL("user32")
	openClipboard    = user32.MustFindProc("OpenClipboard")
	closeClipboard   = user32.MustFindProc("CloseClipboard")
	emptyClipboard   = user32.MustFindProc("EmptyClipboard")
	getClipboardData = user32.MustFindProc("GetClipboardData")
	setClipboardData = user32.MustFindProc("SetClipboardData")

	kernel32     = syscall.NewLazyDLL("kernel32")
	globalAlloc  = kernel32.NewProc("GlobalAlloc")
	globalFree   = kernel32.NewProc("GlobalFree")
	globalLock   = kernel32.NewProc("GlobalLock")
	globalUnlock = kernel32.NewProc("GlobalUnlock")
	lstrcpy      = kernel32.NewProc("lstrcpyW")
)

func readAll() (string, error) {
	r, _, err := openClipboard.Call(0)
	if r == 0 {
		return "", err
	}
	defer closeClipboard.Call()

	h, _, err := getClipboardData.Call(cfUnicodetext)
	if r == 0 {
		return "", err
	}

	l, _, err := globalLock.Call(h)
	if l == 0 {
		return "", err
	}

	text := syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(l))[:])

	r, _, err = globalUnlock.Call(h)
	if r == 0 {
		return "", err
	}

	return text, nil
}

func writeAll(text string) error {
	r, _, err := openClipboard.Call(0)
	if r == 0 {
		return err
	}
	defer closeClipboard.Call()

	r, _, err = emptyClipboard.Call(0)
	if r == 0 {
		return err
	}

	data := syscall.StringToUTF16(text)

	h, _, err := globalAlloc.Call(gmemFixed, uintptr(len(data)*int(unsafe.Sizeof(data[0]))))
	if h == 0 {
		return err
	}

	l, _, err := globalLock.Call(h)
	if l == 0 {
		return err
	}

	r, _, err = lstrcpy.Call(l, uintptr(unsafe.Pointer(&data[0])))
	if r == 0 {
		return err
	}

	r, _, err = globalUnlock.Call(h)
	if r == 0 {
		return err
	}

	r, _, err = setClipboardData.Call(cfUnicodetext, h)
	if r == 0 {
		return err
	}
	return nil
}
