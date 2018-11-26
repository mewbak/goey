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

	filter := (*uint16)(nil)
	if len(m.filters) > 0 {
		buffer := make([]uint16, 0, 1024)
		for _, v := range m.filters {
			tmp, _ := syscall.UTF16FromString(v.name)
			buffer = append(buffer, tmp...)
			tmp, _ = syscall.UTF16FromString(v.pattern)
			buffer = append(buffer, tmp...)
		}
		buffer = append(buffer, 0, 0)
		filter = &buffer[0]
	}

	file := [1024]uint16{}

	ofn := win.OPENFILENAME{
		LStructSize: uint32(unsafe.Sizeof(win.OPENFILENAME{})),
		HwndOwner:   win.HWND(m.handle),
		LpstrFilter: filter,
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
