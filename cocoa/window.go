package cocoa

/*
#include "cocoa.h"
*/
import "C"
import "unsafe"

type Window struct {
private int
}

func NewWindow(title string, width, height uint) *Window {
	println("newWindow")
	ctitle := C.CString(title)
	handle := C.windowNew(ctitle, C.unsigned(width), C.unsigned(height))
	return (*Window)(handle)
}

func (w *Window) Close() {
	C.windowClose(unsafe.Pointer(w))
}

func (w *Window) Uintptr() uintptr {
	return uintptr(unsafe.Pointer(w))
}
