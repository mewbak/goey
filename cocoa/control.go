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
