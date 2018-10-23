package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Decoration is a wrapper for a GDecoration.
type Decoration struct {
	View
	private int
}

func NewDecoration(window *View) *Decoration {
	handle := C.decorationNew(unsafe.Pointer(window))
	return (*Decoration)(handle)
}

func (w *Decoration) Close() {
	C.viewClose(unsafe.Pointer(w))
}
