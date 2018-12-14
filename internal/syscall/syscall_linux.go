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
// static void goey_set_bounds(GtkWidget* handle, gint x, gint y, gint width, gint height ) {
//	GtkWidget* parent = gtk_widget_get_parent( handle );
//	GtkLayout* layout = GTK_LAYOUT(parent);
//	gtk_layout_move( layout, handle, x, y );
//	GtkAllocation alloc = { x, y, width, height };
//	gtk_widget_size_allocate( handle, &alloc );
//}
// static void goey_set_key_info(void* widget, void* evt_, guint r, gchar release ) {
//  GdkEventKey* evt = evt_;
//  evt->type = release ? GDK_KEY_RELEASE : GDK_KEY_PRESS;
//  evt->window = gtk_widget_get_window( GTK_WIDGET(widget) );
//  evt->time = GDK_CURRENT_TIME;
//  evt->send_event = 1;
//  switch ( r ) {
//    case 0x1b:
//      evt->keyval = GDK_KEY_Escape;
//      evt->hardware_keycode = 9;
//      break;
//    case '\n':
//      evt->keyval = GDK_KEY_Return;
//      evt->hardware_keycode = 36;
//      break;
//    default: evt->keyval = r;
//  }
// }
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

// SetBounds is a specialized wrapper around gtk_widget_size_allocate.  However,
// this function also assumes that the parent is a GtkLayout, and so also
// moves the widget using gtk_layout_move.
func SetBounds(widget *gtk.Widget, x, y, width, height int) {
	p := unsafe.Pointer(widget.GObject)
	C.goey_set_bounds((*C.GtkWidget)(p), C.gint(x), C.gint(y), C.gint(width), C.gint(height))
}

func SetEventKeyInformation(widget *gtk.Widget, evt *gdk.EventKey, r rune, release uint8) {
	C.goey_set_key_info(unsafe.Pointer(widget.Native()), unsafe.Pointer(evt.GdkEvent), C.guint(r), C.gchar(release))
}
