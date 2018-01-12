package goey

import (
	"github.com/lxn/win"
	"image"
)

func (w *Empty) mount(parent NativeWidget) (MountedWidget, error) {
	retval := &mountedEmpty{}

	return retval, nil
}

type mountedEmpty struct {
}

func (w *mountedEmpty) Close() {
	// Virtual control, so no resources to release
}

func (w *mountedEmpty) MeasureWidth() (DIP, DIP) {
	return 13, 13
}

func (w *mountedEmpty) MeasureHeight(width DIP) (DIP, DIP) {
	// Same as static text
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 13, 13
}

func (w *mountedEmpty) SetBounds(bounds image.Rectangle) {
	// Virtual control, so no resource to resize
}

func (w *mountedEmpty) SetOrder(hwnd win.HWND) win.HWND {
	return hwnd
}
