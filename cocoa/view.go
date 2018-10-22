package cocoa

/*
#include "cocoa.h"
*/
import "C"
import "unsafe"

// View is a wrapper for a NSView.
type View struct {
	private int
}

// SetFrame is a wrapper for the setFrame message.
func (c *View) SetFrame(x, y, dx, dy int) {
	C.viewSetFrame(unsafe.Pointer(c), C.int(x), C.int(y), C.int(dx), C.int(dy))
}
