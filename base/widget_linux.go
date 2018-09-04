package base

import "github.com/gotk3/gotk3/gtk"

// Control is an opaque type used as a platform-specific handle to a control
// created using the platform GUI.  As an example, this will refer to a HWND
// when targeting Windows, but a *GtkContainer when targeting GTK.
//
// Unless developping new widgets, users should not need to use this type.
//
// Any methods on this type will be platform specific.
type Control struct {
	Handle *gtk.Container
}

// NativeElement contains platform-specific methods that all widgets
// must support on GTK.
type NativeElement interface {
}
