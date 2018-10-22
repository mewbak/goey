package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Button is a wrapper for a NSButton.
type TextField struct {
	Control
	private int
}

func NewTextField(window *Window, title string) *TextField {
	ctitle := C.CString(title)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	handle := C.textfieldNew(unsafe.Pointer(window), ctitle)
	return (*TextField)(handle)
}
