package cocoa

/*
#include "cocoa.h"
*/
import "C"
import "unsafe"

type Button struct {
	Control
	private int
}

func NewButton(window *Window, title string) *Button {
	println("newWindow")
	ctitle := C.CString(title)
	handle := C.buttonNew(unsafe.Pointer(window), ctitle)
	return (*Button)(handle)
}

func (w *Button) Close() {
	C.buttonClose(unsafe.Pointer(w))
}

func (w *Button) Uintptr() uintptr {
	return uintptr(unsafe.Pointer(w))
}
