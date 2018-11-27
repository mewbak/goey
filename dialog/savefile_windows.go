package dialog

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

func (m *SaveFile) show() (string, error) {
	title, err := syscall.UTF16PtrFromString(m.title)
	if err != nil {
		return "", err
	}

	file := [1024]uint16{}

	ofn := win.OPENFILENAME{
		LStructSize: uint32(unsafe.Sizeof(win.OPENFILENAME{})),
		HwndOwner:   win.HWND(m.handle),
		LpstrFilter: buildFilterString(m.filters),
		LpstrFile:   &file[0],
		NMaxFile:    1024,
		LpstrTitle:  title,
	}

	rc := win.GetSaveFileName(&ofn)
	if !rc {
		if err := win.CommDlgExtendedError(); err != 0 {
			return "", fmt.Errorf("call to GetOpenFileName failed with code %x", err)
		}
		return "", nil
	}
	return syscall.UTF16ToString(file[:]), nil
}

// WithOwner sets the owner of the dialog box.
func (m *SaveFile) WithOwner(hwnd win.HWND) *SaveFile {
	m.handle = uintptr(hwnd)
	return m
}
