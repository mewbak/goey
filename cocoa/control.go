package cocoa

/*
#include "cocoa.h"
*/
import "C"
import "unsafe"

type Control struct {
	View
	private int
}

func toBOOL(b bool) C.BOOL {
	if b {
		return C.YES
	}
	return C.NO
}

func (c *Control) Close() {
	C.controlClose(unsafe.Pointer(c))
}

func (c *Control) IsEnabled() bool {
	rc := C.controlIsEnabled(unsafe.Pointer(c))
	return rc != 0
}

func (c *Control) IntrinsicContentSize() (int, int) {
	var h C.int
	w := C.controlIntrinsicContentSize(unsafe.Pointer(c), &h)
	return int(w), int(h)
}

func (c *Control) SetEnabled(enabled bool) {
	C.controlSetEnabled(unsafe.Pointer(c), toBOOL(enabled))
}
