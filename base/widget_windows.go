package base

import (
	"github.com/lxn/win"
)

// Control is an opaque type used as a platform-specific handle to a control
// created using the platform GUI.  As an example, this will refer to a HWND
// when targeting Windows, but a *GtkWidget when targeting GTK.
//
// Unless developping new widgets, users should not need to use this type.
//
// Any method's on this type will be platform specific.
type Control struct {
	HWnd win.HWND
}

// NativeElement contains platform-specific methods that all widgets
// must support on WIN32
type NativeElement interface {
	SetOrder(previous win.HWND) win.HWND
}
