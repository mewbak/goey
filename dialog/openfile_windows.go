package dialog

import (
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

func (m *OpenFile) show() (string, error) {
	title, err := syscall.UTF16PtrFromString(m.title)
	if err != nil {
		return "", err
	}

	file := [1024]uint16{}

	ofn := win.OPENFILENAME{
		LStructSize: uint32(unsafe.Sizeof(win.OPENFILENAME{})),
		HwndOwner:   win.HWND(m.handle),
		LpstrFile:   &file[0],
		NMaxFile:    1024,
		LpstrTitle:  title,
	}

	rc := win.GetOpenFileName(&ofn)
	if !rc {
		return "", nil
	}
	return syscall.UTF16ToString(file[:]), nil
}

// WithOwner sets the owner of the dialog box.
func (m *OpenFile) WithOwner(hwnd win.HWND) *OpenFile {
	m.handle = uintptr(hwnd)
	return m
}
