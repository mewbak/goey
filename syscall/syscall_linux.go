package syscall

import (
	"unsafe"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
)

// #cgo pkg-config: gdk-3.0
// #include <gdk/gdk.h>
import "C"

func fromBool(value bool) C.gboolean {
	if value {
		return C.TRUE
	}
	return C.FALSE
}

func PixbufNewFromBytes(bytes []uint8, colorspace gdk.Colorspace, hasAlpha bool, bitsPerSample int,
	width int, height int, rowStride int) *gdk.Pixbuf {

	ret := C.gdk_pixbuf_new_from_data(
		(*C.guchar)(&bytes[0]), C.GdkColorspace(colorspace), fromBool(hasAlpha), C.int(bitsPerSample),
		C.int(width), C.int(height), C.int(rowStride), nil, nil)
	return &gdk.Pixbuf{glib.Take(unsafe.Pointer(ret))}
}
