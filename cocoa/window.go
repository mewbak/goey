package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

type Window struct {
	private int
}

func NewWindow(title string, width, height uint) *Window {
	ctitle := C.CString(title)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	handle := C.windowNew(ctitle, C.unsigned(width), C.unsigned(height))
	return (*Window)(handle)
}

func (w *Window) Close() {
	C.windowClose(unsafe.Pointer(w))
}

func (w *Window) ContentSize() (int, int) {
	var h C.int
	px := C.windowContentSize(unsafe.Pointer(w), &h)
	return int(px), int(h)
}
