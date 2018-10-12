package cocoa

/*
#include "cocoa.h"
*/
import "C"
import "unsafe"

type View struct {
	private int
}

func (c *View) SetFrame(x, y, dx, dy int) {
	C.viewSetFrame(unsafe.Pointer(c), C.int(x), C.int(y), C.int(dx), C.int(dy))
}
