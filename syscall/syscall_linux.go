// Package syscall fills in some missing APIs from GTK+ 3 that are not provided
// by gotk3's package.  These are limited to those required by the package goey,
// and should be candidates for upstreaming.  Naming convention is to convert
// snake-case used by the C API to camel-case.
//
// This package is intended for internal use.
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
// static void makeRectangle(void* rc_, int x, int y, int w, int h ) {
//	GdkRectangle *rc = rc_;
//	rc->x = x; rc->y = y; rc->width = w; rc->height = h;
// }
// static void goey_set_bounds(GtkWidget* handle, gint x, gint y, gint width, gint height ) {
//	GtkWidget* parent = gtk_widget_get_parent( handle );
//	GtkLayout* layout = GTK_LAYOUT(parent);
//	gtk_layout_move( layout, handle, x, y );
//	GtkAllocation alloc = { x, y, width, height };
//	gtk_widget_size_allocate( handle, &alloc );
//}
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

// LayoutGetVAdjustment is a wrapper around gtk_layout_get_vadjustment.
func LayoutGetVAdjustment(widget *gtk.Layout) *gtk.Adjustment {
	p := unsafe.Pointer(widget.GObject)
	a := C.gtk_layout_get_vadjustment((*C.GtkLayout)(p))
	obj := glib.Take(unsafe.Pointer(a))
	return &gtk.Adjustment{glib.InitiallyUnowned{obj}}
}

// WidgetGetAllocation is a wrapper around gtk_widget_get_allocation.
func WidgetGetAllocation(widget *gtk.Widget) (int, int, int, int) {
	var a C.GtkAllocation
	p := unsafe.Pointer(widget.GObject)
	C.gtk_widget_get_allocation((*C.GtkWidget)(p), &a)
	return int(a.x), int(a.y), int(a.width), int(a.height)
}

// WidgetGetPreferredHeightForWidth is a wrapper around gtk_widget_get_preferred_height_for_width.
func WidgetGetPreferredHeightForWidth(widget *gtk.Widget, width int) (int, int) {
	var minimum, natural C.gint
	p := unsafe.Pointer(widget.GObject)
	C.gtk_widget_get_preferred_height_for_width((*C.GtkWidget)(p), C.gint(width), &minimum, &natural)
	return int(minimum), int(natural)
}

// WindowSetInteractiveDebugging is a wrapper around gtk_window_set_interactive_debugging.
func WindowSetInteractiveDebugging(enable bool) {
	C.gtk_window_set_interactive_debugging(fromBool(enable))
}

// SetBounds is a specialized wrapper around gtk_widget_size_allocate.  However,
// this function also assumes that the parent is a GtkLayout, and so also
// moves the widget using gtk_layout_move.
func SetBounds(widget *gtk.Widget, x, y, width, height int) {
	p := unsafe.Pointer(widget.GObject)
	C.goey_set_bounds((*C.GtkWidget)(p), C.gint(x), C.gint(y), C.gint(width), C.gint(height))
}
