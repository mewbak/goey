package dialog

import (
	"fmt"
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
		LpstrFilter: buildFilterString(m.filters),
		LpstrFile:   &file[0],
		NMaxFile:    1024,
		LpstrTitle:  title,
	}

	rc := win.GetOpenFileName(&ofn)
	if !rc {
		if err := win.CommDlgExtendedError(); err != 0 {
			return "", fmt.Errorf("call to GetOpenFileName failed with code %x", err)
		}
		return "", nil
	}
	return syscall.UTF16ToString(file[:]), nil
}

func buildFilterString(filters []filter) *uint16 {
	// If there are no filters, we want to return a nil pointer.
	// This will let windows select appropriate default behaviour.
	if len(filters) == 0 {
		return nil
	}

	// See documentation for OPENFILENAME structures, but we build a single
	// buffer with pairs of null-terminated strings.
	buffer := make([]uint16, 0, 1024)
	for _, v := range filters {
		tmp, _ := syscall.UTF16FromString(v.name)
		buffer = append(buffer, tmp...)
		tmp, _ = syscall.UTF16FromString(v.pattern)
		buffer = append(buffer, tmp...)
	}

	// Final double null-terminated string marks end of buffer.
	buffer = append(buffer, 0, 0)
	return &buffer[0]
}

// WithOwner sets the owner of the dialog box.
func (m *OpenFile) WithOwner(hwnd win.HWND) *OpenFile {
	m.handle = uintptr(hwnd)
	return m
}
