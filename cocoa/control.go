package cocoa

/*
#include "cocoa.h"
*/
import "C"
import "unsafe"

type Control struct {
	private int
}

func toBOOL(b bool) C.BOOL {
	if b {
		return C.YES
	}
	return C.NO
}

func (c *Control) SetEnabled(enabled bool) {
	C.controlSetEnabled(unsafe.Pointer(c), toBOOL(enabled))
}

func (c *Control) SetBounds(x, y, dx, dy int) {
	C.controlSetBounds(unsafe.Pointer(c), C.int(x), C.int(y), C.int(dx), C.int(dy))
}
