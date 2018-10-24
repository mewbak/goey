package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import (
	"image/color"
	"unsafe"
)

// Decoration is a wrapper for a GDecoration.
type Decoration struct {
	View
	private int
}

func toColor(clr color.Color) C.nscolor_t {
	r, g, b, a := clr.RGBA()
	return C.nscolor_t{
		r: C.uint8_t(r >> 8),
		g: C.uint8_t(g >> 8),
		b: C.uint8_t(b >> 8),
		a: C.uint8_t(a >> 8),
	}
}

func NewDecoration(window *View, fill color.Color, stroke color.Color) *Decoration {
	handle := C.decorationNew(unsafe.Pointer(window), toColor(fill), toColor(stroke))
	return (*Decoration)(handle)
}

func (w *Decoration) Close() {
	C.viewClose(unsafe.Pointer(w))
}

func (w *Decoration) SetBorderRadius(x, y int) {
	radius := C.nssize_t{
		width:  C.int32_t(x),
		height: C.int32_t(y),
	}
	C.decorationSetBorderRadius(unsafe.Pointer(w), radius)
}

func (w *Decoration) SetFillColor(fill color.Color) {
	C.decorationSetFillColor(unsafe.Pointer(w), toColor(fill))
}

func (w *Decoration) SetStrokeColor(fill color.Color) {
	C.decorationSetStrokeColor(unsafe.Pointer(w), toColor(fill))
}
