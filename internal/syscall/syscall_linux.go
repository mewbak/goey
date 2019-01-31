// Package syscall provides platform-dependent routines required to support the
// package goey.
// In particular, when using GTK+3, the goal is to fill in some missing APIs
// that are not provided by gotk3's package.
// Anything found herein should be a candidate for upstreaming.
// Naming convention is to convert snake-case used by the C API to camel-case.
//
// This package is intended for internal use.
//
// This package contains platform-specific details.
package syscall

import (
	"unsafe"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// #cgo pkg-config: gdk-3.0 gtk+-3.0
// #include <gdk/gdk.h>
// #include <gtk/gtk.h>
// #include "syscall_linux.h"
import "C"

func fromBool(value bool) C.gboolean {
	if value {
		return C.TRUE
	}
	return C.FALSE
}

// PixbufNewFromBytes is a wrapper around gdk_pixbuf_new_from_data.
func PixbufNewFromBytes(bytes []uint8, colorspace gdk.Colorspace, hasAlpha bool, bitsPerSample int,
	width int, height int, rowStride int) *gdk.Pixbuf {

	ret := C.gdk_pixbuf_new_from_data(
		(*C.guchar)(&bytes[0]), C.GdkColorspace(colorspace), fromBool(hasAlpha), C.int(bitsPerSample),
		C.int(width), C.int(height), C.int(rowStride), nil, nil)
	return &gdk.Pixbuf{glib.Take(unsafe.Pointer(ret))}
}

// PixbufGetFromWindow is a wrapper around gdk_pixbuf_get_from_window.
func PixbufGetFromWindow(root *gdk.Window, window *gtk.Window) *gdk.Pixbuf {
	// Get the coordinates for the window
	tmp := C.gtk_widget_get_window((*C.GtkWidget)(unsafe.Pointer(window.GObject)))
	var x, y, w, h C.gint
	C.gdk_window_get_origin(tmp, &x, &y)
	C.gdk_window_get_geometry(tmp, nil, nil, &w, &h)

	// The offsets to the dimensions below are to capture the title bar and
	// the borders for the window.  This is tuned to XFCE, and will likely need
	// to be adjusted with any other DE.
	ret := C.gdk_pixbuf_get_from_window((*C.GdkWindow)(unsafe.Pointer(root.Native())),
		x-1, y-25, C.gint(w)+2, C.gint(h)+26)
	if ret == nil {
		return nil
	}
	return &gdk.Pixbuf{glib.Take(unsafe.Pointer(ret))}
}

// WidgetGetPreferredHeightForWidth is a wrapper around gtk_widget_get_preferred_height_for_width.
func WidgetGetPreferredHeightForWidth(widget *gtk.Widget, width int) (int, int) {
	var minimum, natural C.gint
	p := unsafe.Pointer(widget.GObject)
	C.gtk_widget_get_preferred_height_for_width((*C.GtkWidget)(p), C.gint(width), &minimum, &natural)
	return int(minimum), int(natural)
}

// WidgetSendKey is a wrapper around gtk_widget_event to send a key press and release event.
func WidgetSendKey(widget *gtk.Widget, keyval rune, modifiers gdk.ModifierType, release uint8) {
	p := unsafe.Pointer(widget.GObject)
	//C.gtk_test_widget_send_key((*C.GtkWidget)(p), C.guint(keyval), C.GdkModifierType(modifiers))
	C.goey_widget_send_key((*C.GtkWidget)(p), C.guint(keyval), C.GdkModifierType(modifiers), C.gchar(release))
}

// SetBounds is a specialized wrapper around gtk_widget_size_allocate.  However,
// this function also assumes that the parent is a GtkLayout, and so also
// moves the widget using gtk_layout_move.
func SetBounds(widget *gtk.Widget, x, y, width, height int) {
	p := unsafe.Pointer(widget.GObject)
	C.goey_set_bounds((*C.GtkWidget)(p), C.gint(x), C.gint(y), C.gint(width), C.gint(height))
}

var (
	invokeFunction = make(chan func(), 1)
)

// MainContextInvoke is a wrapper around g_main_context.
func MainContextInvoke(function func()) {
	invokeFunction <- function
	C.goey_main_context_invoke()
}

//export mainContextInvokeCallback
func mainContextInvokeCallback() {
	fn := <-invokeFunction
	fn()
}
