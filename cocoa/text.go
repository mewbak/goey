package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Text is a wrapper for a NSText.
type Text struct {
	View
	private int
}

func NewText(window *View, title string) *Text {
	ctitle := C.CString(title)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	handle := C.textNew(unsafe.Pointer(window), ctitle)
	return (*Text)(handle)
}

func (w *Text) SetAlignment(align int) {
	C.textSetAlignment(unsafe.Pointer(w), C.int(align))
}

func (w *Text) SetText(text string) {
	ctext := C.CString(text)
	defer func() {
		C.free(unsafe.Pointer(ctext))
	}()

	C.textSetText(unsafe.Pointer(w), ctext)
}
