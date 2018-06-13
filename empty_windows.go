package goey

import (
	"github.com/lxn/win"
)

func (w *Empty) mount(parent NativeWidget) (Element, error) {
	retval := &mountedEmpty{}

	return retval, nil
}

type mountedEmpty struct {
}

func (w *mountedEmpty) Close() {
	// Virtual control, so no resources to release
}

func (w *mountedEmpty) SetBounds(bounds Rectangle) {
	// Virtual control, so no resource to resize
}

func (w *mountedEmpty) SetOrder(hwnd win.HWND) win.HWND {
	return hwnd
}
